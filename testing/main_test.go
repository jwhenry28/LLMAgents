package main

import (
	"github.com/gocolly/colly/v2"
	"os"
	"path/filepath"
	"testing"
)

// MockWebsite represents your scraper class
type MockWebsite struct {
	collector *colly.Collector
}

// TestData represents the data you want to extract
type TestData struct {
	Title   string
	Content string
}

func NewMockWebsite() *MockWebsite {
	c := colly.NewCollector()
	return &MockWebsite{collector: c}
}

// localFile creates a proper file:// URL from a local path
func localFile(path string) string {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	return "file://" + absPath
}

func TestScraper(t *testing.T) {
	// Create test HTML file
	testHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1 class="title">Hello World</h1>
    <div class="content">Test Content</div>
</body>
</html>
`
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test-*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up after test

	// Write test HTML to file
	if _, err := tmpFile.Write([]byte(testHTML)); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Initialize scraper
	scraper := NewMockWebsite()
	result := &TestData{}

	// Set up collectors
	scraper.collector.OnHTML("h1.title", func(e *colly.HTMLElement) {
		result.Title = e.Text
	})
	scraper.collector.OnHTML("div.content", func(e *colly.HTMLElement) {
		result.Content = e.Text
	})

	// Visit local file
	err = scraper.collector.Visit(localFile(tmpFile.Name()))
	if err != nil {
		t.Fatalf("Failed to visit local file: %v", err)
	}

	// Assert results
	expectedTitle := "Hello World"
	if result.Title != expectedTitle {
		t.Errorf("Expected title %s, got %s", expectedTitle, result.Title)
	}

	expectedContent := "Test Content"
	if result.Content != expectedContent {
		t.Errorf("Expected content %s, got %s", expectedContent, result.Content)
	}
}