package utils

import (
	"log/slog"

	"github.com/gocolly/colly"
	"hackandpray.com/media-curator/model"
)

type Scraper struct {
	URL       string
	Anchors   []model.Anchor
	InnerText string

	collector *colly.Collector
}

func NewScraper(url string) *Scraper {
	collector := colly.NewCollector()
	s := Scraper{
		URL:       url,
		Anchors:   []model.Anchor{},
		InnerText: "",
		collector: collector,
	}
	s.initialize()
	return &s
}

func (s *Scraper) Scrape() {
	s.collector.Visit(s.URL)
	s.collector.Wait()
}

func (s *Scraper) initialize() {
	s.collector.OnRequest(func(r *colly.Request) {
		slog.Debug("visiting: " + s.URL)
	})
	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		s.Anchors = append(s.Anchors, model.NewAnchor(e.Text, e.Attr("href")))
	})
	s.collector.OnHTML("p,article,code,h1,h2,h3,h4,h5,h6", func(e *colly.HTMLElement) {
		s.InnerText += e.Text
	})
	s.collector.OnResponse(func(r *colly.Response) {
		slog.Debug("completed: " + s.URL)
	})
}
