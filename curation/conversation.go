package curation

import (
	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
)

type Conversation struct {
	llm      llm.LLM
	ongoing  bool
	messages []model.Chat
}

func NewConversation() Conversation {
	return Conversation{}
}

func (c *Conversation) initialize() {
	c.llm = llm.NewMockLLM()
	c.messages = []model.Chat{}
}
