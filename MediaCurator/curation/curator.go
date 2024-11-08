package curation

import (
	"fmt"
	"log/slog"
	"strings"

	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
	"hackandpray.com/media-curator/tools"
	"hackandpray.com/media-curator/utils"
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

	messages := c.initialMessages(scraper)
	conversation := NewConversation(c.llm, messages)
	conversation.RunConversation(seed)
}

func (c *Curator) initialMessages(scraper *utils.Scraper) []model.Chat {
	return []model.Chat{
		{
			Role: "system",
			Content: fmt.Sprintf(
				utils.SYSTEM_PROMPT,
				c.getDescription(),
				tools.NewDecide(model.ToolInput{}).Help(),
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
