package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"hackernews_scraper/model"
)

// HackerNewsScraper is responsible for scraping Hacker News.
type HackerNewsScraper struct {
	URL       *url.URL
	Anchors   []model.Anchor
	InnerText string
	collector *colly.Collector
}

// NewHackerNewsScraper initializes a new HackerNewsScraper.
func NewHackerNewsScraper(rawURL string) *HackerNewsScraper {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	return &HackerNewsScraper{
		URL:       parsedURL,
		collector: colly.NewCollector(),
	}
}

// Scrape performs the scraping of the Hacker News page.
func (h *HackerNewsScraper) Scrape() {
	h.collector.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		text := e.Text
		h.Anchors = append(h.Anchors, model.Anchor{Text: text, HRef: href})
	})

	h.collector.OnResponse(func(r *colly.Response) {
		h.InnerText = string(r.Body)
	})

	h.collector.Visit(h.URL.String())
}

// GetFeaturedArticles filters and returns the featured articles.
func (h *HackerNewsScraper) GetFeaturedArticles() []model.Anchor {
	var featured []model.Anchor
	for _, anchor := range h.Anchors {
		if strings.HasPrefix(anchor.HRef, "http") && anchor.Text != "" {
			featured = append(featured, anchor)
		}
	}
	return featured
}

func main() {
	scraper := NewHackerNewsScraper("https://news.ycombinator.com/")
	scraper.Scrape()
	featuredArticles := scraper.GetFeaturedArticles()

	for _, article := range featuredArticles {
		fmt.Printf("Title: %s, URL: %s\n", article.Text, article.HRef)
	}
}