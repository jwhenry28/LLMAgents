package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	articles := []map[string]string{}

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		title := e.ChildText("td.title > a.storylink")
		url := e.ChildAttr("td.title > a.storylink", "href")
		if title != "" && url != "" {
			articles = append(articles, map[string]string{"title": title, "url": url})
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
