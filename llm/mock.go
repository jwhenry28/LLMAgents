package llm

import "hackandpray.com/media-curator/model"

type MockLLM struct {
}

func NewMockLLM() *MockLLM {
	return &MockLLM{}
}

func (llm *MockLLM) CompleteChat(messages []model.Chat) (string, error) {
	return `{ "tool": "decide", "args": [ "NOTIFY", "http://example.com" ]}`, nil
}
