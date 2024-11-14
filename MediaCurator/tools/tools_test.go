package tools

import (
	"strings"
	"testing"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

func genTextToolInput(input string) model.ToolInput {
	toolInput, _ := model.NewTextToolInput(input)
	return toolInput
}

func TestDecide(t *testing.T) {
	testCases := []struct {
		name     string
		input    model.ToolInput
		expected string
		match    bool
	}{
		{
			name:     "notify",
			input:    genTextToolInput("decide NOTIFY https://example.com 'generic justification'"),
			expected: "notified",
			match:    true,
		},
		{
			name:     "ignore",
			input:    genTextToolInput("decide IGNORE https://example.com 'generic justification'"),
			expected: "ignored",
			match:    true,
		},
		{
			name:     "unknown decision",
			input:    genTextToolInput("decide FOOBAR https://example.com 'generic justification'"),
			expected: "unknown decision",
			match:    false,
		},
		{
			name:     "missing justification",
			input:    genTextToolInput("decide NOTIFY https://example.com"),
			expected: "notified",
			match:    false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			tool := NewDecide(test.input)

			if test.match != tool.Match() {
				t.Errorf("Expected Match to return %t but got %t", test.match, tool.Match())
			}

			got := tool.Invoke()
			if got != test.expected {
				t.Errorf("Expected %q to contain %q", got, test.expected)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	testCases := []struct {
		name     string
		input    model.ToolInput
		expected string
		match    bool
	}{
		{
			name:     "valid url",
			input:    genTextToolInput("fetch https://example.com"),
			expected: "Example Domain",
			match:    true,
		},
		{
			name:     "invalid url",
			input:    genTextToolInput("fetch not-a-url"),
			expected: "",
			match:    false,
		},
		{
			name:     "missing args",
			input:    genTextToolInput("fetch"),
			expected: "",
			match:    false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			tool := NewFetch(test.input)

			if test.match != tool.Match() {
				t.Errorf("Expected Match to return %t but got %t", test.match, tool.Match())
			}
			if !test.match {
				return
			}

			actual := tool.Invoke()
			if !strings.Contains(actual, test.expected) {
				t.Errorf("Expected %q to contain %q", actual, test.expected)
			}
		})
	}
}
