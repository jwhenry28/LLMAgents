package scrapers

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

const (
	DEFAULT_FORMATTED_LEN = 5000
)

type DefaultScraper struct {
	*BaseScraper
}

func NewDefaultScraper(urlString string) (Scraper, error) {
	baseScraper, err := NewBaseScraper(urlString)
	if err != nil {
		return nil, err
	}
	s := DefaultScraper{
		BaseScraper: &baseScraper,
	}
	s.initialize()
	return &s, nil
}

func (s *DefaultScraper) initialize() {
	s.Collector.OnRequest(func(r *colly.Request) {
		s.Err = nil
	})
	s.Collector.OnError(func(r *colly.Response, err error) {
		s.StatusCode = r.StatusCode
		s.Err = err
	})
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
	})
	s.Collector.OnHTML("p,article,code,h1,h2,h3,h4,h5,h6", func(e *colly.HTMLElement) {
		s.InnerText += e.Text
		s.FullText += e.Text
	})
	s.Collector.OnScraped(func(r *colly.Response) {
		s.StatusCode = r.StatusCode
	})
}

func (s *DefaultScraper) GetFormattedText() string {
	truncated := s.FullText
	if len(truncated) > DEFAULT_FORMATTED_LEN {
		truncated = truncated[0:DEFAULT_FORMATTED_LEN]
	}
	return truncated
}
