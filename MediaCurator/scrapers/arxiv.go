package scrapers

import (
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
)

type ArxivScraper struct {
	BaseScraper
	Done   bool
	Titles []string
	HRefs  []string
}

func NewArxivScraper(urlString string) (Scraper, error) {
	baseScraper, err := NewBaseScraper(urlString)
	if err != nil {
		return nil, err
	}
	s := ArxivScraper{
		BaseScraper: baseScraper,
		Done:        false,
		Titles:      []string{},
		HRefs:       []string{},
	}
	s.initialize()
	return &s, nil
}

func (s *ArxivScraper) initialize() {
	s.collector.OnHTML("dl", func(e *colly.HTMLElement) {
		if s.isOlderThanToday(e) {
			s.collector.OnHTMLDetach("dl")
			s.Done = true
			return
		}

		e.ForEach("dt", func(i int, e *colly.HTMLElement) {
			abstractLink := e.ChildAttr("a[href^='/abs']", "href")
			if abstractLink != "" {
				abstractUrl := "https://" + s.GetHostname() + abstractLink
				s.HRefs = append(s.HRefs, abstractUrl)
			}
		})

		e.ForEach("dd", func(i int, e *colly.HTMLElement) {
			title := e.ChildText(".list-title")
			title = strings.TrimPrefix(title, "Title:")
			s.Titles = append(s.Titles, strings.TrimSpace(title))
		})
	})

	s.collector.OnError(func(r *colly.Response, err error) {
		slog.Warn("scraper error", "error", err)
		s.StatusCode = r.StatusCode
		s.Err = err
	})
}

func (s *ArxivScraper) isOlderThanToday(e *colly.HTMLElement) bool {
	dateStr := e.ChildText("h3")
	dateStr = strings.Split(dateStr, "(")[0]
	dateStr = strings.TrimSpace(dateStr)

	date, err := time.Parse("Mon, 2 Jan 2006", dateStr)
	if err != nil {
		slog.Warn("failed to parse date", "error", err, "text", dateStr)
		return false
	}

	today := time.Now().YearDay()
	cursor := date.YearDay()

	return today > cursor || (today == 1 && cursor >= 365)
}

func (s *ArxivScraper) Scrape() {
	s.FullText = ""
	s.InnerText = ""
	s.Anchors = []model.Anchor{}

	i := 0
	pageSize := 100
	guardrail := 1000
	for !s.Done {
		s.url.RawQuery = url.Values{
			"skip": {strconv.Itoa(i)},
			"show": {strconv.Itoa(pageSize)},
		}.Encode()

		s.collector.Visit(s.url.String())
		s.collector.Wait()

		i += pageSize
		if i > guardrail {
			slog.Warn("guardrail hit - ignoring results", "guardrail", guardrail)
			return
		}
	}

	if len(s.HRefs) != len(s.Titles) {
		slog.Warn("scraper mismatch", "href count", len(s.HRefs), "title count", len(s.Titles))
	}

	for i, href := range s.HRefs {
		s.Anchors = append(s.Anchors, model.NewAnchor(s.Titles[i], href))
	}
}

func (s *ArxivScraper) GetFormattedText() string {
	formatted := ""
	for _, anchor := range s.GetAnchors() {
		formatted += "Title: " + anchor.Text + "\n"
		formatted += "HRef: " + anchor.HRef + "\n\n"
	}
	return formatted
}
