package scrapers

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

func getTestDataDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd + "/testdata/", nil
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name           string
		constructor    func(string) (Scraper, error)
		url            string
		anchors        []string
		contentSnippet string
	}{
		{
			"DefaultScraper",
			NewDefaultScraper,
			"example.html",
			[]string{"https://www.iana.org/domains/example"},
			"This domain is for use in illustrative examples in documents.",
		},
		{
			"HackerNewsScraper",
			NewHackerNewsScraper,
			"hackernews.html",
			[]string{
				"https://www.johndcook.com/blog/2008/09/19/writes-large-correct-programs/",
				"https://www.nature.com/articles/d41586-024-03756-w",
				"https://docs.maxxinteractive.com/",
				"https://iximiuz.com/en/series/computer-networking-fundamentals/",
				"https://github.com/ColleagueRiley/RGFW",
			},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileUrl := "file:///" + test.url
			scraper, err := test.constructor(fileUrl)
			if err != nil {
				t.Fatalf("failed to construct scraper: %s", test.name)
			}

			testData, err := getTestDataDir()
			if err != nil {
				t.Fatal("failed to retrieve test data directory")
			}

			scraper.SetTransport(http.NewFileTransport(http.Dir(testData)))
			scraper.Scrape()
			for _, expected := range test.anchors {
				found := false
				for _, actual := range scraper.GetAnchors() {
					if expected == actual.HRef {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("failed to retrieve anchor: %s", expected)
				}
			}

			if !strings.Contains(scraper.GetFormattedText(), test.contentSnippet) {
				t.Errorf("failed to retrieve expected content.")
			}
		})
	}
}
