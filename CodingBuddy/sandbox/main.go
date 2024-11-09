package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

func main() {
	c := colly.NewCollector()

	var articles []map[string]string

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if url != "newcomments" && url != "jobs" && url != "submit" {
			if e.DOM.Parent().Is("td.title > a") {
				article := map[string]string{
					"title": e.Text,
					"url":   e.Request.AbsoluteURL(url),
				}
				articles = append(articles, article)
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		jsonData, err := json.MarshalIndent(articles, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))
	})

	err := c.Visit("https://news.ycombinator.com")
	if err != nil {
		log.Fatal(err)
	}

	c.Wait()
	fmt.Println(articles)
}
