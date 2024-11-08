package tools

import (
	"testing"

	"hackandpray.com/llm-agents/model"
)

func TestDecide(t *testing.T) {
	testCases := []struct {
		name     string
		input    model.ToolInput
		expected string
		match    bool
	}{
		{
			name: "notify",
			input: model.ToolInput{
				Name: "decide",
				Args: []string{"NOTIFY", "https://example.com"},
			},
			expected: "notified",
			match:    true,
		},
		{
			name: "ignore",
			input: model.ToolInput{
				Name: "decide",
				Args: []string{"IGNORE", "https://example.com"},
			},
			expected: "ignored",
			match:    true,
		},
		{
			name: "unknown decision",
			input: model.ToolInput{
				Name: "decide",
				Args: []string{"FOOBAR", "https://example.com"},
			},
			expected: "unknown decision",
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
				t.Errorf("Expected %q but got %q", test.expected, got)
			}
		})
	}
}
