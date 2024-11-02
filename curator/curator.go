package curator

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
	fm       utils.FileManager
	seeds    []string
	scrapers map[string]utils.Scraper
	llm      llm.LLM
	ongoing  bool
}

func New() *Curator {
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

	c.llm = llm.NewMockLLM()
}

func (c *Curator) Curate() {
	for _, seed := range c.seeds {
		c.process(seed)
	}
}

func (c *Curator) process(seed string) {
	c.scrapeSeed(seed)
	c.runLLMSession(seed)
}

func (c *Curator) scrapeSeed(seed string) {
	scraper := c.getOrCreateScraper(seed)
	scraper.Scrape()
}

func (c *Curator) runLLMSession(seed string) {
	c.ongoing = true
	messages := c.initializeMessages(seed)
	for c.ongoing {
		responseJSON, err := c.llm.CompleteChat(messages)
		if err != nil {
			slog.Error("LLM session failed", "err", err)
			return
		}
		messages = append(messages, model.NewChat("assistant", responseJSON))

		output := c.runTool(responseJSON)
		userMsg := model.NewChat("user", output)
		messages = append(messages, userMsg)
	}
}

func (c *Curator) initializeMessages(seed string) []model.Chat {
	scraper := c.scrapers[seed]
	messages := []model.Chat{
		{
			Role:    "system",
			Content: fmt.Sprintf(utils.SYSTEM_PROMPT, c.getDescription()),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.URL_PROMPT, seed, scraper.InnerText),
		},
	}

	return messages
}

func (c *Curator) getOrCreateScraper(seed string) *utils.Scraper {
	scraper, ok := c.scrapers[seed]
	if !ok {
		scraper = *utils.NewScraper(seed)
	}

	return &scraper
}

func (c *Curator) getDescription() string {
	return c.fm.Read("data/description.txt")
}

func (c *Curator) runTool(input string) string {
	var tool model.Tool
	err := json.Unmarshal([]byte(input), &tool)
	if err != nil {
		return fmt.Sprintf("LLM response is not conforming JSON: %s\nerror: %s\n", tool, err.Error())
	}

	if tool.Name == "decide" {
		args := tool.Args
		if len(args) < 2 {
			return fmt.Sprintf("too few args to 'decide' call, expected 2: %v\n", args)
		}

		if args[0] == "NOTIFY" {
			slog.Info("Sending notification for url", "url", args[1])
			c.ongoing = false
			return ""
		} else if args[0] == "IGNORE" {
			slog.Info("Ignoring url", "url", args[1])
			c.ongoing = false
			return ""
		} else {
			return "improperly formatted 'decide' call: " + args[0]
		}
	}

	return "unknown tool: " + tool.Name
}
