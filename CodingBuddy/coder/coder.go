package coder

import (
	"github.com/jwhenry28/LLMAgents/shared/conversation"
	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Coder struct {
	llm llm.LLM
}

func NewCoder(llm llm.LLM) *Coder {
	return &Coder{llm: llm}
}

func (c *Coder) Code() {
	conversationIsOver := func(c *conversation.Conversation) bool {
		modelResponse := c.Messages[len(c.Messages)-2]
		selectedTool, _ := model.ToolInputFromJSON(modelResponse.Content)

		toolOutput := c.Messages[len(c.Messages)-1]

		return selectedTool.Name == "decide" && (toolOutput.Content == "notified" || toolOutput.Content == "ignored")
	}

	messages := c.initialMessages(scraper)
	conversation := conversation.NewConversation(c.llm, messages, conversationIsOver)
	conversation.RunConversation(seed)
}
