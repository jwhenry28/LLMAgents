package curation

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

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
)

var Registry = map[string]func(string) (scrapers.Scraper, error){
	"news.ycombinator.com": scrapers.NewHackerNewsScraper,
}

type Result struct {
	URL           string
	Decision      string
	Justification string
}

type Curator struct {
	fm       utils.FileManager
	seeds    []string
	scrapers map[string]scrapers.Scraper
	llm      llm.LLM
}

func NewCurator(llm llm.LLM) *Curator {
	c := Curator{
		fm:  utils.NewFileManager(),
		llm: llm,
	}

	c.registerTools()
	c.loadSeeds()
	c.loadScrapers()
	return &c
}

func (c *Curator) registerTools() {
	tools.RegisterTool("help", tools.NewHelp)
	tools.Registry["fetch"] = local.NewFetch
	tools.Registry["decide"] = local.NewDecide
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
	c.scrapers = make(map[string]scrapers.Scraper)
	for _, seed := range c.seeds {
		scraper, err := c.getOrCreateScraper(seed)
		if err != nil {
			slog.Warn("Error creating seed scraper", "error", err)
			continue
		}
		c.scrapers[seed] = scraper
	}
}

func (c *Curator) Curate() {
	for _, seed := range c.seeds {
		c.runLLMSession(seed)
	}
}

func (c *Curator) getOrCreateScraper(seed string) (scrapers.Scraper, error) {
	seed = strings.TrimSpace(seed)
	seed = strings.ToLower(seed)
	seed = strings.TrimSuffix(seed, "/")
	seed = strings.TrimPrefix(seed, "https://")
	seed = strings.TrimPrefix(seed, "http://")
	constructor, ok := Registry[seed]
	if !ok {
		constructor = scrapers.NewDefaultScraper
	}
	scraper, err := constructor(seed)
	if err != nil {
		return nil, fmt.Errorf("error creating seed scraper: %s", err)
	}
	return scraper, nil
}

func (c *Curator) runLLMSession(seed string) {
	conversationIsOver := func(c conversation.Conversation) bool {
		messages := c.GetMessages()
		modelResponse := messages[len(messages)-1]
		selectedTool, _ := model.NewJSONToolInput(modelResponse.Content)

		constructor, ok := tools.Registry["decide"]

		return ok && selectedTool.GetName() == "decide" && constructor(selectedTool).Match()
	}

	scraper := c.scrapeSeed(seed)
	var decision model.ToolInput
	results := []Result{}
	anchors := scraper.GetAnchors()
	for _, anchor := range anchors {
		subScraper, err := c.getOrCreateScraper(anchor.HRef)
		if err != nil {
			slog.Warn("Error getting sub-scraper", "error", err)
			continue
		}
		subScraper.Scrape()
		llmMessages := c.initialMessages(subScraper)
		conversation := conversation.NewChatConversation(c.llm, llmMessages, conversationIsOver, "json")
		conversation.RunConversation()
		messages := conversation.GetMessages()
		lastMessage := messages[len(messages)-1]
		decision, err = model.NewJSONToolInput(lastMessage.Content)
		if err != nil {
			slog.Warn("Error parsing decision", "error", err)
			continue
		}
		args := decision.GetArgs()
		result := Result{
			Decision:      args[0],
			URL:           args[1],
			Justification: args[2],
		}
		results = append(results, result)
		c.saveResults(results)
	}
}

func (c *Curator) saveResults(results []Result) {
	resultsJson, err := json.Marshal(results)
	if err != nil {
		slog.Error("Error marshaling results", "error", err)
		return
	}

	err = os.WriteFile("./data/results.json", resultsJson, 0644)
	if err != nil {
		slog.Error("Error writing results file", "error", err)
		return
	}
}

func (c *Curator) scrapeSeed(seed string) scrapers.Scraper {
	scraper, err := c.getOrCreateScraper(seed)
	if err != nil {
		return nil
	}
	scraper.Scrape()
	return scraper
}

func (c *Curator) initialMessages(scraper scrapers.Scraper) []model.Chat {
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
			Content: fmt.Sprintf(utils.CONTENT_PROMPT, scraper.GetURL(), scraper.GetFormattedAnchors(), scraper.GetFormattedText()),
		},
	}
}

func (c *Curator) getDescription() string {
	return c.fm.Read(DESCRIPTION_FILE)
}
