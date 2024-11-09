package curation

import (
	"fmt"
	"log/slog"
	"strings"

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

type Curator struct {
	fm       utils.FileManager
	seeds    []string
	scrapers map[string]*utils.Scraper
	llm      llm.LLM
}

func NewCurator(llm llm.LLM) *Curator {
	c := Curator{
		fm:  utils.NewFileManager(),
		llm: llm,
	}

	c.loadSeeds()
	c.loadScrapers()
	return &c
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
	c.scrapers = make(map[string]*utils.Scraper)
	for _, seed := range c.seeds {
		scraper, err := utils.NewScraper(seed)
		if err != nil {
			slog.Warn("Error creating seed scraper", "error", err)
			continue
		}
		c.scrapers[seed] = scraper
	}
}

func (c *Curator) Curate() {
	for _, seed := range c.seeds {
		c.scrapeSeed(seed)
		c.runLLMSession(seed)
	}
}

func (c *Curator) scrapeSeed(seed string) {
	scraper, err := c.getOrCreateScraper(seed)
	if err != nil {
		return
	}
	scraper.Scrape()
}

func (c *Curator) getOrCreateScraper(seed string) (*utils.Scraper, error) {
	scraper, ok := c.scrapers[seed]
	var err error = nil
	if !ok {
		scraper, err = utils.NewScraper(seed)
		if err == nil {
			c.scrapers[seed] = scraper
		}
	}

	return scraper, err
}

func (c *Curator) runLLMSession(seed string) {
	scraper, err := c.getOrCreateScraper(seed)
	if err != nil {
		slog.Error("Error getting scraper", "error", err)
		return
	}

	conversationIsOver := func(c *conversation.Conversation) bool {
		modelResponse := c.Messages[len(c.Messages)-2]
		selectedTool, _ := model.ToolInputFromJSON(modelResponse.Content)
	
		toolOutput := c.Messages[len(c.Messages)-1]
	
		return selectedTool.Name == "decide" && (toolOutput.Content == "notified" || toolOutput.Content == "ignored")
	}

	messages := c.initialMessages(scraper)
	conversation := conversation.NewConversation(c.llm, messages, conversationIsOver)
	conversation.RunConversation()
}

func (c *Curator) initialMessages(scraper *utils.Scraper) []model.Chat {
	return []model.Chat{
		{
			Role: "system",
			Content: fmt.Sprintf(
				utils.SYSTEM_PROMPT,
				c.getDescription(),
				local.NewDecide(model.ToolInput{}).Help(),
				tools.NewHelp(model.ToolInput{Name: "help", Args: []string{}}).Invoke(),
			),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.CONTENT_PROMPT, scraper.URL, scraper.GetAnchorString(), scraper.InnerText),
		},
	}
}

func (c *Curator) getDescription() string {
	return c.fm.Read(DESCRIPTION_FILE)
}
