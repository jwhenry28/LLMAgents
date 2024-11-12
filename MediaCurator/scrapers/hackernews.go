package scrapers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func Old() {
	c := colly.NewCollector()

	var articles []map[string]string

	c.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if !strings.Contains(url, "vote") && !strings.Contains(url, "hide") && !strings.Contains(url, "item") {
			title := e.Text
			if !strings.HasPrefix(url, "http") {
				url = "https://news.ycombinator.com/" + url
			}
			article := map[string]string{"title": title, "url": url}
			articles = append(articles, article)
		}
	})

	err := c.Visit("https://news.ycombinator.com")
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
