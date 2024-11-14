package scrapers

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type DefaultScraper struct {
	BaseScraper
}

func NewDefaultScraper(urlString string) (Scraper, error) {
	baseScraper, err := NewBaseScraper(urlString)
	if err != nil {
		return nil, err
	}
	s := DefaultScraper{
		BaseScraper: baseScraper,
	}
	s.initialize()
	return &s, nil
}

func (s *DefaultScraper) initialize() {
	s.Collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		hyperlink := e.Attr("href")
		if !strings.HasPrefix(hyperlink, "http") {
			url := s.GetURL()
			if !strings.HasSuffix(url, "/") && !strings.HasPrefix(hyperlink, "/") {
				url += "/"
			}
			hyperlink = url + hyperlink
		}
		anchor := model.NewAnchor(e.Text, hyperlink)
		s.Anchors = append(s.Anchors, anchor)
		s.FullText += fmt.Sprintf("<a href=\"%s\">%s</a>", hyperlink, e.Text)
	})
	s.Collector.OnHTML("p,article,code,h1,h2,h3,h4,h5,h6", func(e *colly.HTMLElement) {
		s.InnerText += e.Text
		s.FullText += e.Text
	})
}
