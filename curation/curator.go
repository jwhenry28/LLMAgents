package curation

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
	"hackandpray.com/media-curator/utils"
)

type Curator struct {
	fm           utils.FileManager
	seeds        []string
	scrapers     map[string]utils.Scraper
	conversation Conversation
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
	c.ongoing = true
	c.initializeMessages(seed)
	for c.ongoing {
		c.generateModelResponse()
		c.runTool()
	}
}

func (c *Curator) initializeMessages(seed string) {
	scraper := c.scrapers[seed]
	c.messages = []model.Chat{
		{
			Role:    "system",
			Content: fmt.Sprintf(utils.SYSTEM_PROMPT, c.getDescription()),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.URL_PROMPT, seed, scraper.InnerText),
		},
	}
}

func (c *Curator) generateModelResponse() {
	raw, err := c.llm.CompleteChat(c.messages)
	var response model.Chat
	if err != nil {
		slog.Error("LLM session failed", "err", err)
		c.ongoing = false
		response = model.NewChat("", "")
	} else {
		response = model.NewChat("assistant", raw)
	}

	c.messages = append(c.messages, response)
}

func (c *Curator) getDescription() string {
	return c.fm.Read("data/description.txt")
}

func (c *Curator) runTool() {
	input := c.messages[len(c.messages)-1].Content
	var tool model.Tool
	err := json.Unmarshal([]byte(input), &tool)

	output := "unknown tool: " + tool.Name
	if err != nil {
		output = fmt.Sprintf("LLM response is not conforming JSON: %s\nerror: %s\n", tool, err.Error())
	}

	if tool.Name == "decide" {
		args := tool.Args
		if len(args) < 2 {
			output = fmt.Sprintf("too few args to 'decide' call, expected 2: %v\n", args)
		}

		if args[0] == "NOTIFY" {
			slog.Info("Sending notification for url", "url", args[1])
			c.ongoing = false
			output = ""
		} else if args[0] == "IGNORE" {
			slog.Info("Ignoring url", "url", args[1])
			c.ongoing = false
			output = ""
		} else {
			output = "improperly formatted 'decide' call: " + args[0]
		}
	}

	c.messages = append(c.messages, model.NewChat("user", output))
}
