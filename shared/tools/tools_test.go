package tools

import (
	"testing"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

func TestGetUsage(t *testing.T) {
	tests := []struct {
		name     string
		toolType string
		input    string
		expected string
	}{
		{"JSONToolInput", model.JSON_TOOL_TYPE, `{ "tool": "help", "args": [ "help" ]}`, `usage: { "tool": "help", "args": [ <tool-name> ]}`},
		{"TextToolInput", model.TEXT_TOOL_TYPE, `help help`, `usage: help <tool-name>`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			helpInput, err := model.NewToolInput(test.toolType, test.input)
			if err != nil {
				t.Errorf("failed to create tool input: %v", err)
			}

			help := NewHelp(helpInput)

			if help.Usage() != test.expected {
				t.Errorf("incorrect usage;\ngot      %q\nexpected %q", help.Usage(), test.expected)
			}
		})
	}
}
