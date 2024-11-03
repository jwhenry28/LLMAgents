package llm

import (
	"hackandpray.com/media-curator/model"
)

type MockLLM struct {
	messages []model.Chat
}

func NewMockLLM() *MockLLM {
	return &MockLLM{}
}

func (llm *MockLLM) AddMessage(message model.Chat) {
	llm.messages = append(llm.messages, message)
}

func (llm *MockLLM) CompleteChat(_ []model.Chat) (string, error) {
	if len(llm.messages) == 0 {
		return "no messages available", nil
	}

	message := llm.messages[0]
	llm.messages = llm.messages[1:]
	return message.Content, nil
}
