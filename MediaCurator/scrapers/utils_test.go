package scrapers

import "testing"

func TestScraperConstructor(t *testing.T) {
	testCases := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{"http://example.com/", "http://example.com", false},
		{"http://example.com", "http://example.com", false},
		{"example.com", "https://example.com", false},
		{"example", "", true},
	}

	for _, test := range testCases {
		_, err := NewScraper(test.input)
		if test.expectErr == (err == nil) {
			t.Errorf("Expected error for input %q but got nil", test.input)
		}
	}
}
