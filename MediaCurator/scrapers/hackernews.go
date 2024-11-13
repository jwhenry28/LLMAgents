package scrapers

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type HackerNewsScraper struct {
	BaseScraper
}

func NewHackerNewsScraper(urlString string) (Scraper, error) {
	baseScraper, err := NewBaseScraper(urlString)
	if err != nil {
		return nil, err
	}
	s := HackerNewsScraper{
		BaseScraper: baseScraper,
	}
	s.initialize()
	return &s, nil
}

func (s *HackerNewsScraper) initialize() {
	s.Collector.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if !strings.Contains(url, "vote") && !strings.Contains(url, "hide") && !strings.Contains(url, "item") {
			title := e.Text
			if !strings.HasPrefix(url, "http") {
				url = s.GetURL() + url
			}
			s.Anchors = append(s.Anchors, model.NewAnchor(title, url))
		}
	})
}
