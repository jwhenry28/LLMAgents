package scrapers

import (
	"log/slog"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
	"golang.org/x/net/publicsuffix"
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
	s.Collector.OnRequest(func(r *colly.Request) {
		s.Err = nil
	})
	s.Collector.OnError(func(r *colly.Response, err error) {
		slog.Warn("scraper error", "error", err)
		s.StatusCode = r.StatusCode
		s.Err = err
	})
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
	hostRoot, _ := publicsuffix.EffectiveTLDPlusOne(s.URL.Hostname())
	targetRoot, err := publicsuffix.EffectiveTLDPlusOne(url.Hostname())
	return err == nil && hostRoot != targetRoot
}

func (s *HackerNewsScraper) GetFormattedText() string {
	formatted := ""
	for _, anchor := range s.GetAnchors() {
		formatted += "Title: " + anchor.Text + "\n"
		formatted += "HRef: " + anchor.HRef + "\n\n"
	}

	return formatted
}
