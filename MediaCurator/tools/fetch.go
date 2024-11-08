package tools

import (
	"net/url"

	"hackandpray.com/llm-agents/model"
	"hackandpray.com/llm-agents/utils"
)

type Fetch struct {
	Base
}

func NewFetch(input model.ToolInput) Tool {
	brief := "fetch: fetches the content of the specified URL."
	usage := `usage: { "tool": "fetch", "args": [ <url> ]}
args:
- url: The URL you wish to fetch content from. Must start with http or https.`

	return Fetch{Base: Base{Input: input, BriefText: brief, UsageText: usage}}
}

func (task Fetch) Match() bool {
	if len(task.Input.Args) < 1 {
		return false
	}

	_, err := url.ParseRequestURI(task.Input.Args[0])
	return err == nil
}

// TODO: use common data store with curator.scrapers
func (task Fetch) Invoke() string {
	scraper, err := utils.NewScraper(task.Input.Args[0])
	if err != nil {
		return "error: " + err.Error()
	}
	
	scraper.Scrape()
	return scraper.InnerText
}
