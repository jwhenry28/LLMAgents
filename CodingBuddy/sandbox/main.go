package sandbox

import (
	"strings"
	"testing"

	"github.com/gocolly/colly"
)

// Article represents a HackerNews article
type Article struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// ScrapeArticles is the callback function that extracts article data
func ScrapeArticles(e *colly.HTMLElement, articles *[]Article) {
	// HN articles are in a table with class="itemlist"
	// Each article row has class="athing"
	if e.Attr("class") != "athing" {
		return
	}

	titleElement := e.ChildText("td.title > span.titleline > a")
	url := e.ChildAttr("td.title > span.titleline > a", "href")

	// Skip internal links and non-article links
	if url == "" || isInternalLink(url) {
		return
	}

	*articles = append(*articles, Article{
		Title: titleElement,
		URL:   url,
	})
}

// isInternalLink checks if the URL is an internal HN link
func isInternalLink(url string) bool {
	internalPaths := []string{
		"/jobs",
		"/newcomments",
		"/submit",
		"/item",
		"/user",
		"/newest",
		"/front",
		"/best",
	}

	for _, path := range internalPaths {
		if strings.HasPrefix(url, path) {
			return true
		}
	}
	return false
}

// TestScrapeArticles contains all the test cases
func TestScrapeArticles(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected []Article
	}{
		{
			name: "valid article",
			html: `
				<tr class="athing">
					<td class="title">
						<span class="titleline">
							<a href="https://example.com">Test Article</a>
						</span>
					</td>
				</tr>`,
			expected: []Article{
				{Title: "Test Article", URL: "https://example.com"},
			},
		},
		{
			name: "internal link should be ignored",
			html: `
				<tr class="athing">
					<td class="title">
						<span class="titleline">
							<a href="/jobs">Jobs</a>
						</span>
					</td>
				</tr>`,
			expected: []Article{},
		},
		{
			name: "multiple articles",
			html: `
				<tr class="athing">
					<td class="title">
						<span class="titleline">
							<a href="https://example1.com">Article 1</a>
						</span>
					</td>
				</tr>
				<tr class="athing">
					<td class="title">
						<span class="titleline">
							<a href="https://example2.com">Article 2</a>
						</span>
					</td>
				</tr>`,
			expected: []Article{
				{Title: "Article 1", URL: "https://example1.com"},
				{Title: "Article 2", URL: "https://example2.com"},
			},
		},
		{
			name: "non-article row should be ignored",
			html: `
				<tr class="something-else">
					<td class="title">
						<span class="titleline">
							<a href="https://example.com">Not an Article</a>
						</span>
					</td>
				</tr>`,
			expected: []Article{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new collector for each test
			c := colly.NewCollector()
			var articles []Article

			// Register the callback
			c.OnHTML("tr", func(e *colly.HTMLElement) {
				ScrapeArticles(e, &articles)
			})

			// Visit the HTML string
			c.OnHTML("html", func(e *colly.HTMLElement) {
				e.DOM.Find("tr").Each(func(_ int, s *colly.HTMLElement) {
					ScrapeArticles(s, &articles)
				})
			})

			err := c.Visit("data:text/html," + tt.html)
			if err != nil {
				t.Fatalf("Failed to visit HTML: %v", err)
			}

			// Compare results
			if len(articles) != len(tt.expected) {
				t.Errorf("Expected %d articles, got %d", len(tt.expected), len(articles))
				return
			}

			for i, article := range articles {
				if article != tt.expected[i] {
					t.Errorf("Article %d: expected %+v, got %+v", i, tt.expected[i], article)
				}
			}
		})
	}
}

// TestIsInternalLink tests the internal link detection
func TestIsInternalLink(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com", false},
		{"/jobs", true},
		{"/newcomments", true},
		{"/submit", true},
		{"/item?id=123", true},
		{"/user?id=test", true},
		{"http://news.ycombinator.com", false},
		{"/newest", true},
		{"/front", true},
		{"/best", true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := isInternalLink(tt.url)
			if result != tt.expected {
				t.Errorf("isInternalLink(%q) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}
