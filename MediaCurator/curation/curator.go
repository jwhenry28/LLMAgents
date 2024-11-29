package curation

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	localmodel "github.com/jwhenry28/LLMAgents/media-curator/model"
	"github.com/jwhenry28/LLMAgents/media-curator/scrapers"
	local "github.com/jwhenry28/LLMAgents/media-curator/tools"
	"github.com/jwhenry28/LLMAgents/media-curator/utils"
	"github.com/jwhenry28/LLMAgents/shared/conversation"
	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

const (
	SEEDS_FILE       = "data/seeds.txt"
	DESCRIPTION_FILE = "data/description.txt"
	TOOL_TYPE        = model.JSON_TOOL_TYPE
)

type ScraperConstructor func(string) (scrapers.Scraper, error)

var ScrapersRegistry = map[string]ScraperConstructor{
	"news.ycombinator.com": scrapers.NewHackerNewsScraper,
	"arxiv.org":            scrapers.NewArxivScraper,
	"nautil.us":            scrapers.NewNautilusScraper,
	"medium.com":           scrapers.NewMediumScraper,
}

type Curator struct {
	fm        utils.FileManager
	seeds     []string
	decisions []local.Decision
	scrapers  map[string]scrapers.Scraper
	profile   *localmodel.Profile
	llm       llm.LLM
	sendEmail bool
	filename  string
}

func NewCurator(llm llm.LLM, profile *localmodel.Profile) *Curator {
	c := Curator{
		fm:        utils.NewFileManager(),
		llm:       llm,
		decisions: []local.Decision{},
		profile:   profile,
	}

	c.registerTools()
	c.loadSeeds()
	c.loadScrapers()

	return &c
}

func (c *Curator) SetSendEmail(sendEmail bool) {
	c.sendEmail = sendEmail
}

func (c *Curator) registerTools() {
	tools.RegisterTool("help", tools.NewHelp)
	tools.RegisterTool("fetch", local.NewFetch)
	tools.RegisterTool("complete", local.NewComplete)
}

func (c *Curator) loadSeeds() {
	lines := strings.Split(c.fm.Read(SEEDS_FILE), "\n")
	seeds := []string{}
	for _, url := range lines {
		if strings.TrimSpace(url) != "" {
			seeds = append(seeds, url)
		}
	}

	c.seeds = seeds
}

func (c *Curator) loadScrapers() {
	if len(c.seeds) == 0 {
		slog.Warn("loading curator scrapers without any seeds")
	}

	c.scrapers = make(map[string]scrapers.Scraper)
	for _, seed := range c.seeds {
		scraper, err := c.getOrCreateScraper(seed)
		if err != nil {
			slog.Warn("error creating seed scraper", "error", err)
			continue
		}
		c.scrapers[seed] = scraper
	}
}

func (c *Curator) getOrCreateScraper(seed string) (scrapers.Scraper, error) {
	constructor := c.getScraperConstructor(seed)
	scraper, err := constructor(seed)
	if err != nil {
		return nil, err
	}
	return scraper, nil
}

func (c *Curator) getScraperConstructor(seed string) ScraperConstructor {
	constructor, ok := ScrapersRegistry[c.formatSeed(seed)]
	if !ok {
		constructor = scrapers.NewDefaultScraper
	}
	return constructor
}

func (c *Curator) formatSeed(seed string) string {
	u, err := url.Parse(seed)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func (c *Curator) GetRecipientEmail() string {
	return c.profile.Email
}

func (c *Curator) GetRecipientName() string {
	return c.profile.Name
}

func (c *Curator) Curate() {
	slog.Info("curating for profile", "email", c.GetRecipientEmail(), "name", c.GetRecipientName())
	for _, seed := range c.seeds {
		c.runLLMSession(seed)
	}

	if c.GetRecipientEmail() != "" && c.sendEmail {
		c.sendResultsEmail()
	}
}

func (c *Curator) runLLMSession(seed string) {
	scraper := c.scrapeSeed(seed)
	slog.Info("processing seed", "seed", seed, "anchors", len(scraper.GetAnchors()))

	messages, err := c.generateInitialMessages(scraper)
	if err != nil {
		slog.Warn("error getting sub-scraper", "error", err)
		return
	}

	lastMsg, err := c.generateModelDecisions(messages)
	if err != nil {
		slog.Error("error retrieving decisions", "error", err)
		return
	}

	decisions, err := model.NewToolInput(TOOL_TYPE, lastMsg.Content)
	if err != nil {
		slog.Error("error parsing decision", "error", err)
		return
	}

	c.processDecision(decisions)
	slog.Info("completed seed", "seed", seed)
}

func (c *Curator) processDecision(decisions model.ToolInput) {
	for _, arg := range decisions.GetArgs() {
		var decision local.Decision
		err := json.Unmarshal([]byte(arg), &decision)
		if err != nil {
			slog.Warn("error unmarshaling decision", "decision", arg, "error", err)
			continue
		}

		c.decisions = append(c.decisions, decision)
	}
	c.saveResults()
}

func (c *Curator) extractFinalTool(conversation conversation.Conversation) (model.ToolInput, error) {
	lastMessage := conversation.GetLastMessage()
	return model.NewJSONToolInput(lastMessage.Content)
}

func (c *Curator) generateModelDecisions(messages []model.Chat) (model.Chat, error) {
	conversationIsOver := func(conv conversation.Conversation) bool {
		finalTool, err := c.extractFinalTool(conv)
		if err != nil {
			return false
		}

		endToolName := "complete"
		decideConstructor, ok := tools.Registry[endToolName]
		return ok && finalTool.GetName() == endToolName && decideConstructor(finalTool).Invoke() == "completed"
	}

	conversation := conversation.NewChatConversation(c.llm, messages, conversationIsOver, TOOL_TYPE, true)
	conversation.RunConversation()
	return conversation.GetLastMessage(), nil
}

func (c *Curator) saveResults() {
	results := map[string]interface{}{
		"recipient":   c.GetRecipientEmail(),
		"description": c.GetDescription(),
		"seeds":       c.seeds,
		"decisions":   c.decisions,
	}
	resultsJson, err := json.Marshal(results)
	if err != nil {
		slog.Error("Error marshaling results", "error", err)
		return
	}

	dataFolder := fmt.Sprintf("./data/%s", time.Now().Format(time.DateOnly))
	err = os.MkdirAll(dataFolder, 0755)
	if err != nil {
		slog.Error("Error creating data directory", "error", err)
		return
	}

	if c.filename == "" {
		c.filename = fmt.Sprintf("%s_%s_%s.json", c.llm.Type(), c.GetRecipientName(), time.Now().Format(time.TimeOnly))
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dataFolder, c.filename), resultsJson, 0644)
	if err != nil {
		slog.Error("Error writing results file", "error", err)
		return
	}
}

func (c *Curator) sendResultsEmail() {
	mailer, err := utils.NewEmailSender(c.llm.Type() + "@hackandpray.com")
	if err != nil {
		slog.Error("Error creating email sender", "error", err)
		return
	}

	body := c.buildEmail()
	err = mailer.SendEmail(c.GetRecipientEmail(), "Media Curator Results", body)
	if err != nil {
		slog.Error("Error sending email", "error", err)
	}
}

func (c *Curator) buildEmail() string {
	seedsBlob := strings.Join(c.seeds, "\n")
	articles := c.getPickedArticles()

	articlesBlob := "Unfortunately, I didn't find any articles I think would interest you today."
	if len(articles) > 0 {
		articlesBlob = "I've curated the following articles for you to read:\n"
		articlesBlob += strings.Join(articles, "\n")
	}

	return fmt.Sprintf(utils.EMAIL_TEMPLATE, seedsBlob, articlesBlob, c.llm.Type())
}

func (c *Curator) getPickedArticles() []string {
	articles := []string{}
	for _, result := range c.decisions {
		if result.Decision == "NOTIFY" {
			articles = append(articles, fmt.Sprintf("<a href=\"%s\">%s</a>", result.URL, result.Title))
		}
	}

	return articles
}

func (c *Curator) scrapeSeed(seed string) scrapers.Scraper {
	scraper, err := c.getOrCreateScraper(seed)
	if err != nil {
		return nil
	}
	scraper.Scrape()
	return scraper
}

func (c *Curator) generateInitialMessages(scraper scrapers.Scraper) ([]model.Chat, error) {
	return c.formatInitialMessages(scraper), nil

}

func (c *Curator) formatInitialMessages(scraper scrapers.Scraper) []model.Chat {
	dummy, _ := model.NewToolInput(TOOL_TYPE, "")
	return []model.Chat{
		{
			Role: "system",
			Content: fmt.Sprintf(
				utils.SYSTEM_PROMPT,
				c.GetDescription(),
				c.getToolList(),
				local.NewComplete(dummy).Help(),
				model.JSON_FORMAT_MSG, //TODO: encapsulate this with TOOL_TYPE
			),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.CONTENT_PROMPT, scraper.GetURL(), scraper.GetFormattedText()),
		},
	}
}

func (c *Curator) getToolList() string {
	return tools.NewHelp(model.JSONToolInput{Name: "help", Args: []string{}}).Invoke()
}

func (c *Curator) GetDescription() string {
	return c.profile.Interests
}
