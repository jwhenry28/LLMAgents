package curation

import (
	"fmt"
	"strings"

	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
	"hackandpray.com/media-curator/utils"
)

type Curator struct {
	fm       utils.FileManager
	seeds    []string
	scrapers map[string]utils.Scraper
}

func NewCurator() *Curator {
	c := Curator{}
	c.initialize()
	return &c
}

func (c *Curator) initialize() {
	c.fm = utils.NewFileManager()
	c.seeds = strings.Split(c.fm.Read("data/seeds.txt"), "\n")
	c.scrapers = make(map[string]utils.Scraper)

	for _, url := range c.seeds {
		c.scrapers[url] = *utils.NewScraper(url)
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
		scraper = *utils.NewScraper(seed)
	}

	return &scraper
}

func (c *Curator) runLLMSession(seed string) {
	scraper := c.getOrCreateScraper(seed)
	conversation := NewConversation(llm.NewHuman(), c.initialMessages(seed, scraper.InnerText))
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
			Content: fmt.Sprintf(utils.URL_PROMPT, seed, seedHTML),
		},
	}
}

func (c *Curator) getDescription() string {
	return c.fm.Read("data/description.txt")
}
