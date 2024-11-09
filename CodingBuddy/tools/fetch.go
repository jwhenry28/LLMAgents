package tools

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type scraper struct {
	collector   *colly.Collector
	ScrapedText string
}

type Fetch struct {
	scraper *scraper
	tools.Base
}

func NewFetch(input model.ToolInput) tools.Tool {
	brief := "fetch: issues a GET request to the specified URL and returns the raw contents."
	usage := `usage: { "tool": "fetch", "args": [ <url> ]}
args:
- url: The URL you wish to fetch content from. Must start with http or https.`

	collector := colly.NewCollector()
	scraper := &scraper{
		collector:   collector,
		ScrapedText: "",
	}
	collector.OnHTML("p,article,code,h1,h2,h3,h4,h5,h6", func(e *colly.HTMLElement) {
		scraper.ScrapedText += e.Text + "\n"
	})
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		scraper.ScrapedText += fmt.Sprintf("<a href=\"%s\">%s</a>\n", e.Attr("href"), e.Text)
	})

	return Fetch{scraper: scraper, Base: tools.Base{Input: input, BriefText: brief, UsageText: usage}}
}

func (task Fetch) Match() bool {
	if len(task.Input.Args) < 1 {
		return false
	}

	_, err := url.ParseRequestURI(task.Input.Args[0])
	return err == nil
}

func (task Fetch) Invoke() string {
	scraper := task.scraper
	scraper.collector.Visit(task.Input.Args[0])
	scraper.collector.Wait()

	return scraper.ScrapedText
}
