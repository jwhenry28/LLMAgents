package scrapers

import (
	"net/url"

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
		hyperlink := e.Attr("href")
		if s.isExternalUrl(hyperlink) {
			s.Anchors = append(s.Anchors, model.NewAnchor(e.Text, hyperlink))
		}
	})
}

func (s *HackerNewsScraper) isExternalUrl(urlString string) bool {
	url, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	return url.Host != s.URL.Host
}
