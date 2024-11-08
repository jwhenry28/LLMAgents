package llm

import (
	"testing"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

func TestMockLLM(t *testing.T) {
	mock := NewMockLLM()
	mock.AddMessage(model.Chat{
		Role:    "user",
		Content: "Hello, world!",
	})
	mock.AddMessage(model.Chat{
		Role:    "user",
		Content: "Hello, again!",
	})

	testCases := []struct {
		name     string
		expected string
	}{
		{"first message", "Hello, world!"},
		{"second message", "Hello, again!"},
		{"no message", "no messages available"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, _ := mock.CompleteChat([]model.Chat{})
			if got != test.expected {
				t.Errorf("Expected %q but got %q", test.expected, got)
			}
		})
	}
}
