package tools

import (
	"fmt"
	"net/url"

	"github.com/jwhenry28/LLMAgents/media-curator/scrapers"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Fetch struct {
	tools.Base
}

func NewFetch(input model.ToolInput) tools.Tool {
	name := "fetch"
	args := []string{"url"}
	brief := "fetch: fetches a preview (5000 chars) of the specified URL's content"
	explanation := `args:
- url: The URL you wish to fetch content from. Must start with http or https.`

	return Fetch{Base: tools.Base{Input: input, Name: name, Args: args, BriefText: brief, ExplanationText: explanation}}
}

func (task Fetch) Match() bool {
	if len(task.Input.GetArgs()) < 1 {
		return false
	}

	_, err := url.ParseRequestURI(task.Input.GetArgs()[0])
	return err == nil
}

func (task Fetch) Invoke() string {
	scraper, err := scrapers.NewDefaultScraper(task.Input.GetArgs()[0])
	if err != nil {
		return "error: " + err.Error()
	}

	scraper.Scrape()
	if scraper.GetErr() != nil {
		return fmt.Sprintf("error: %d - %s", scraper.GetStatusCode(), scraper.GetErr().Error())
	}

	return scraper.GetFormattedText()
}
