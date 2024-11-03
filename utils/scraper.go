package utils

import (
	"log/slog"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"hackandpray.com/media-curator/model"
)

type Scraper struct {
	URL       *url.URL
	Anchors   []model.Anchor
	InnerText string

	collector *colly.Collector
}

func NewScraper(urlString string) (*Scraper, error) {
	url, err := url.ParseRequestURI(formatURL(urlString))
	if err != nil {
		return nil, err
	}

	s := Scraper{
		URL:       url,
		Anchors:   []model.Anchor{},
		InnerText: "",
		collector: colly.NewCollector(),
	}
	s.initialize()
	return &s, nil
}

func formatURL(url string) string {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	url = strings.TrimSuffix(url, "/")

	return url
}

func (s *Scraper) initialize() {
	s.collector.OnRequest(func(r *colly.Request) {
		slog.Debug("visiting: " + s.GetURL())
	})
	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		hyperlink := e.Attr("href")
		if !strings.HasPrefix(hyperlink, "http") {
			hyperlink = s.GetURL() + hyperlink
		}
		s.Anchors = append(s.Anchors, model.NewAnchor(e.Text, hyperlink))
	})
	s.collector.OnHTML("p,article,code,h1,h2,h3,h4,h5,h6", func(e *colly.HTMLElement) {
		s.InnerText += e.Text
	})
	s.collector.OnResponse(func(r *colly.Response) {
		slog.Debug("completed: " + s.GetURL())
	})
}

func (s *Scraper) GetURL() string {
	return s.URL.String()
}

func (s *Scraper) Scrape() {
	s.collector.Visit(s.GetURL())
	s.collector.Wait()
}

func (s *Scraper) GetAnchorString() string {
	output := ""
	for _, anchor := range s.Anchors {
		output += " - " + anchor.AsString() + "\n"
	}
	return output
}
