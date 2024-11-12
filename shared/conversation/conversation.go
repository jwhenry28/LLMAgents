package conversation

import (
	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Conversation interface {
	RunConversation()
	GetMessages() []model.Chat
}

type Base struct {
	llm    llm.LLM
	isOver func(Conversation) bool

	Messages []model.Chat
}

func (b *Base) GetMessages() []model.Chat {
	return b.Messages
}
