package scrapers

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type Scraper interface {
	Scrape()
	GetURL() string
	GetAnchors() []model.Anchor
	GetFormattedAnchors() string
	GetFormattedText() string
}

type BaseScraper struct {
	URL       *url.URL
	Anchors   []model.Anchor
	InnerText string

	Collector *colly.Collector
}

func NewBaseScraper(urlString string) (BaseScraper, error) {
	url, err := url.ParseRequestURI(formatURL(urlString))
	if err != nil {
		return BaseScraper{}, err
	}

	return BaseScraper{
		URL:       url,
		Anchors:   []model.Anchor{},
		InnerText: "",
		Collector: colly.NewCollector(),
	}, nil
}

func formatURL(url string) string {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	url = strings.TrimSuffix(url, "/")

	return url
}

func (s *BaseScraper) GetURL() string {
	return s.URL.String()
}

func (s *BaseScraper) Scrape() {
	s.InnerText = ""
	s.Anchors = []model.Anchor{}
	s.Collector.Visit(s.GetURL())
	s.Collector.Wait()
}

func (s *BaseScraper) GetAnchors() []model.Anchor {
	return s.Anchors
}

func (s *BaseScraper) GetFormattedAnchors() string {
	bytes, err := json.Marshal(s.Anchors)
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

func (s *BaseScraper) GetFormattedText() string {
	return s.InnerText
}
