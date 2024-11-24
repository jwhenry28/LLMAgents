package curation

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/mail"
	"net/url"
	"os"
	"strings"
	"time"

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
	TOOL_TYPE        = "json"
)

type ScraperConstructor func(string) (scrapers.Scraper, error)

var ScrapersRegistry = map[string]ScraperConstructor{
	"news.ycombinator.com": scrapers.NewHackerNewsScraper,
}

type Result struct {
	Decision      string
	URL           string
	Justification string
}

type Curator struct {
	fm        utils.FileManager
	seeds     []string
	results   []Result
	scrapers  map[string]scrapers.Scraper
	llm       llm.LLM
	recipient string
	filename  string
}

func NewCurator(llm llm.LLM, recipient string) *Curator {
	_, err := mail.ParseAddress(recipient)
	if recipient != "" && err != nil {
		slog.Warn("invalid email address, ignoring", "address", recipient)
		recipient = ""
	}
	c := Curator{
		fm:        utils.NewFileManager(),
		llm:       llm,
		results:   []Result{},
		recipient: recipient,
	}

	c.registerTools()
	c.loadSeeds()
	c.loadScrapers()
	return &c
}

func (c *Curator) registerTools() {
	tools.RegisterTool("help", tools.NewHelp)
	tools.RegisterTool("fetch", local.NewFetch)
	tools.RegisterTool("decide", local.NewDecide)
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

func (c *Curator) Curate() {
	for _, seed := range c.seeds {
		c.runLLMSession(seed)
	}

	if c.recipient != "" {
		c.sendResultsEmail()
	}
}

func (c *Curator) runLLMSession(seed string) {
	scraper := c.scrapeSeed(seed)
	anchors := scraper.GetAnchors()
	slog.Info("processing seed", "seed", seed, "anchors", len(anchors))

	for _, anchor := range scraper.GetAnchors() {
		messages, err := c.generateInitialMessages(anchor.HRef)
		if err != nil {
			slog.Warn("error getting sub-scraper", "error", err)
			continue
		}

		decision, err := c.generateModelDecision(messages)
		if err != nil {
			slog.Warn("error parsing decision", "error", err)
			continue
		}

		c.processDecision(decision)
	}
	
	slog.Info("completed seed", "seed", seed)
}

func (c *Curator) processDecision(decision model.ToolInput) {
	args := decision.GetArgs()
	c.results = append(c.results, Result{Decision: args[0], URL: args[1], Justification: args[2]})
	c.saveResults()
}

func (c *Curator) extractFinalTool(conversation conversation.Conversation) (model.ToolInput, error) {
	messages := conversation.GetMessages()
	lastMessage := messages[len(messages)-1]
	return model.NewJSONToolInput(lastMessage.Content)
}

func (c *Curator) generateModelDecision(messages []model.Chat) (model.ToolInput, error) {
	conversationIsOver := func(conv conversation.Conversation) bool {
		finalTool, err := c.extractFinalTool(conv)
		if err != nil {
			return false
		}

		decideConstructor, ok := tools.Registry["decide"]
		return ok && finalTool.GetName() == "decide" && decideConstructor(finalTool).Match()
	}

	conversation := conversation.NewChatConversation(c.llm, messages, conversationIsOver, TOOL_TYPE, true)
	conversation.RunConversation()
	return c.extractFinalTool(conversation)
}

func (c *Curator) saveResults() {
	resultsJson, err := json.Marshal(c.results)
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
		c.filename = fmt.Sprintf("%s_%s.json", c.llm.Type(), time.Now().Format(time.TimeOnly))
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dataFolder, c.filename), resultsJson, 0644)
	if err != nil {
		slog.Error("Error writing results file", "error", err)
		return
	}
}

func (c *Curator) sendResultsEmail() {
	mailer, err := utils.NewEmailSender("joseph@hackandpray.com")
	if err != nil {
		slog.Error("Error creating email sender", "error", err)
		return
	}

	body := c.buildEmail()
	mailer.SendEmail("joseph@hackandpray.com", "Media Curator Results", body)
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
	for _, result := range c.results {
		if result.Decision == "NOTIFY" {
			articles = append(articles, result.URL)
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

func (c *Curator) generateInitialMessages(href string) ([]model.Chat, error) {
	subScraper, err := c.getOrCreateScraper(href)
	if err != nil {
		return nil, err
	}

	subScraper.Scrape()
	return c.formatInitialMessages(subScraper), nil

}

func (c *Curator) formatInitialMessages(scraper scrapers.Scraper) []model.Chat {
	return []model.Chat{
		{
			Role: "system",
			Content: fmt.Sprintf(
				utils.SYSTEM_PROMPT,
				c.getDescription(),
				local.NewDecide(model.JSONToolInput{}).Help(),
				tools.NewHelp(model.JSONToolInput{Name: "help", Args: []string{}}).Invoke(),
				utils.JSON_TOOL_FORMAT,
			),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.CONTENT_PROMPT, scraper.GetURL(), scraper.GetFormattedText()),
		},
	}
}

func (c *Curator) getDescription() string {
	return c.fm.Read(DESCRIPTION_FILE)
}
