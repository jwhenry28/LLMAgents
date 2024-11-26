package scrapers

import (
	"log/slog"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type NautilusScraper struct {
	BaseScraper
}

func NewNautilusScraper(urlString string) (Scraper, error) {
	baseScraper, err := NewBaseScraper(urlString)
	if err != nil {
		return nil, err
	}
	s := NautilusScraper{
		BaseScraper: baseScraper,
	}
	s.initialize()
	return &s, nil
}

func (s *NautilusScraper) initialize() {
	seen := make(map[string]bool)
	s.collector.OnRequest(func(r *colly.Request) {
		s.Err = nil
	})
	s.collector.OnError(func(r *colly.Response, err error) {
		slog.Warn("scraper error", "error", err)
		s.StatusCode = r.StatusCode
		s.Err = err
	})
	s.collector.OnHTML("div.article-box", func(e *colly.HTMLElement) {
		title := e.ChildText("h3 a")
		hyperlink := e.ChildAttr("h3 a", "href")

		_, ok := seen[hyperlink]
		if !ok && s.isInternalUrl(hyperlink) {
			s.Anchors = append(s.Anchors, model.NewAnchor(title, hyperlink))
			seen[hyperlink] = true
		}
	})
}

func (s *NautilusScraper) GetFormattedText() string {
	formatted := ""
	for _, anchor := range s.GetAnchors() {
		formatted += "Title: " + anchor.Text + "\n"
		formatted += "HRef: " + anchor.HRef + "\n\n"
	}

	return formatted
}
