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

	Messages         []model.Chat
	InputConstructor func(string) (model.ToolInput, error)
	Verbose bool
}

func (b *Base) GetMessages() []model.Chat {
	return b.Messages
}

func (b *Base) GetLastMessage() model.Chat {
	return b.Messages[len(b.Messages)-1]
}
