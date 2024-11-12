package main

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/gocolly/colly/v2"
    "sandbox/model"
)

// HackerNewsScraper struct
type HackerNewsScraper struct {
    URL       *url.URL
    Anchors   []model.Anchor
    InnerText string
    collector *colly.Collector
}

// NewHackerNewsScraper constructor
func NewHackerNewsScraper(hnURL string) (*HackerNewsScraper, error) {
    parsedURL, err := url.Parse(hnURL)
    if err != nil {
        return nil, err
    }

    return &HackerNewsScraper{
        URL:       parsedURL,
        Anchors:   []model.Anchor{},
        InnerText: "",
        collector: colly.NewCollector(),
    }, nil
}

// Scrape method
func (hns *HackerNewsScraper) Scrape() error {
    hns.collector.OnHTML("a", func(e *colly.HTMLElement) {
        href := e.Attr("href")
        text := e.Text
        anchor := model.Anchor{Text: text, HRef: href}
        hns.Anchors = append(hns.Anchors, anchor)
    })

    hns.collector.OnResponse(func(r *colly.Response) {
        hns.InnerText = string(r.Body)
    })

    return hns.collector.Visit(hns.URL.String())
}

// GetFeaturedArticles method
func (hns *HackerNewsScraper) GetFeaturedArticles() []model.Anchor {
    var featured []model.Anchor
    for _, anchor := range hns.Anchors {
        if strings.HasPrefix(anchor.HRef, "http") && anchor.Text != "" {
            featured = append(featured, anchor)
        }
    }
    return featured
}

func main() {
    hnURL := "https://news.ycombinator.com/"
    scraper, err := NewHackerNewsScraper(hnURL)
    if err != nil {
        fmt.Println("Error creating scraper:", err)
        return
    }

    err = scraper.Scrape()
    if err != nil {
        fmt.Println("Error scraping:", err)
        return
    }

    featuredArticles := scraper.GetFeaturedArticles()
    for _, article := range featuredArticles {
        fmt.Printf("Title: %s, URL: %s\n", article.Text, article.HRef)
    }
}
