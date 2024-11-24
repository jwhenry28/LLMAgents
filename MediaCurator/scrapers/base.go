package scrapers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type Scraper interface {
	Scrape()
	SetTransport(http.RoundTripper)
	GetURL() string
	GetHostname() string
	GetErr() error
	GetStatusCode() int
	GetAnchors() []model.Anchor
	GetFormattedAnchors() string
	GetInnerText() string
	GetFormattedText() string
}

type BaseScraper struct {
	URL        *url.URL
	Anchors    []model.Anchor
	InnerText  string
	FullText   string
	Err        error
	StatusCode int

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

func (s *BaseScraper) SetTransport(transport http.RoundTripper) {
	s.Collector.WithTransport(transport)
}

func formatURL(url string) string {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "file://") {
		url = "https://" + url
	}

	url = strings.TrimSuffix(url, "/")

	return url
}

func (s *BaseScraper) GetURL() string {
	return s.URL.String()
}

func (s *BaseScraper) GetHostname() string {
	return s.URL.Hostname()
}

func (s *BaseScraper) GetErr() error {
	return s.Err
}

func (s *BaseScraper) GetStatusCode() int {
	return s.StatusCode
}

func (s *BaseScraper) Scrape() {
	s.FullText = ""
	s.InnerText = ""
	s.Anchors = []model.Anchor{}
	s.Collector.Visit(s.GetURL())
	s.Collector.Wait()
}

func (s *BaseScraper) GetAnchors() []model.Anchor {
	seen := make(map[string]bool)
	unique := make([]model.Anchor, 0)

	for _, anchor := range s.Anchors {
		if !seen[anchor.HRef] {
			seen[anchor.HRef] = true
			unique = append(unique, anchor)
		}
	}

	return unique
}

func (s *BaseScraper) GetFormattedAnchors() string {
	bytes, err := json.Marshal(s.Anchors)
	if err != nil {
		return "[]"
	}
	return string(bytes)
}

func (s *BaseScraper) GetInnerText() string {
	return s.InnerText
}

func (s *BaseScraper) GetFormattedText() string {
	return s.FullText
}
