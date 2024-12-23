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
	usage := `usage: fetch <url>
args:
- url: The URL you wish to fetch content from. Must start with http or https.`

	collector := colly.NewCollector()
	scraper := &scraper{
		collector:   collector,
		ScrapedText: "",
	}
	collector.OnHTML("p,article,code,h1,h2,h3,h4,h5,h6", func(e *colly.HTMLElement) {
		scraper.ScrapedText += e.Text
	})
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		dom := e.DOM
		attributes := dom.Nodes[0].Attr
		tag := "<a "
		for i, attribute := range attributes {
			tag += fmt.Sprintf("%s=\"%s\"", attribute.Key, attribute.Val)
			if i < len(attributes)-1 {
				tag += " "
			}
		}
		tag += ">" + e.Text + "</a>\n"
		scraper.ScrapedText += tag
	})

	return Fetch{scraper: scraper, Base: tools.Base{Input: input, BriefText: brief, UsageText: usage}}
}

func (task Fetch) Match() bool {
	args := task.Input.GetArgs()
	if len(args) < 1 {
		return false
	}

	_, err := url.ParseRequestURI(args[0])
	return err == nil
}

func (task Fetch) Invoke() string {
	scraper := task.scraper
	args := task.Input.GetArgs()
	scraper.collector.Visit(args[0])
	scraper.collector.Wait()

	return scraper.ScrapedText
}
