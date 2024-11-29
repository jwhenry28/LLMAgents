package scrapers

import (
	"log/slog"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type MediumScraper struct {
	BaseScraper
}

func NewMediumScraper(urlString string) (Scraper, error) {
	baseScraper, err := NewBaseScraper(urlString)
	if err != nil {
		return nil, err
	}
	s := MediumScraper{
		BaseScraper: baseScraper,
	}
	s.initialize()
	return &s, nil
}

func (s *MediumScraper) initialize() {
	s.collector.OnRequest(func(r *colly.Request) {
		s.Err = nil
	})
	s.collector.OnError(func(r *colly.Response, err error) {
		slog.Warn("scraper error", "error", err)
		s.StatusCode = r.StatusCode
		s.Err = err
	})
	s.collector.OnHTML("a[rel=\"noopener follow\"]", func(e *colly.HTMLElement) {
		hyperlink := e.Attr("href")

		title := strings.TrimSpace(e.ChildText("h2"))
		subtitle := strings.TrimSpace(e.ChildText("h3"))
		if subtitle != "" {
			title += ": " + subtitle
		}

		if s.isInternalUrl(hyperlink) && title != "" {
			path := strings.Split(hyperlink, "?")[0]
			hyperlink = "https://medium.com" + path
			s.Anchors = append(s.Anchors, model.NewAnchor(title, hyperlink))
		}
	})
}

func (s *MediumScraper) GetFormattedText() string {
	formatted := ""
	for _, anchor := range s.GetAnchors() {
		formatted += "Title: " + anchor.Text + "\n"
		formatted += "HRef: " + anchor.HRef + "\n\n"
	}

	return formatted
}
