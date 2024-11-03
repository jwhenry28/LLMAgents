package curation

import (
	"fmt"
	"strings"

	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
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
		c.scrapers[seed] = utils.NewScraper(seed)
	}
}

func (c *Curator) Curate() {
	for _, seed := range c.seeds {
		c.scrapeSeed(seed)
		c.runLLMSession(seed)
	}
}

func (c *Curator) scrapeSeed(seed string) {
	scraper := c.getOrCreateScraper(seed)
	scraper.Scrape()
}

func (c *Curator) getOrCreateScraper(seed string) *utils.Scraper {
	scraper, ok := c.scrapers[seed]
	if !ok {
		scraper = utils.NewScraper(seed)
	}

	return scraper
}

func (c *Curator) runLLMSession(seed string) {
	scraper := c.getOrCreateScraper(seed)
	messages := c.initialMessages(seed, scraper.InnerText)
	conversation := NewConversation(c.llm, messages)
	conversation.RunConversation(seed)
}

func (c *Curator) initialMessages(seed, seedHTML string) []model.Chat {
	return []model.Chat{
		{
			Role:    "system",
			Content: fmt.Sprintf(utils.SYSTEM_PROMPT, c.getDescription()),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.CONTENT_PROMPT, seed, seedHTML),
		},
	}
}

func (c *Curator) getDescription() string {
	return c.fm.Read(DESCRIPTION_FILE)
}
