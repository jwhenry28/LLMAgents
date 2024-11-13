package conversation

import (
	"fmt"
	"log/slog"

	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"

	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type ChatConversation struct {
	Base
}

func NewChatConversation(convoModel llm.LLM, initMessages []model.Chat, isOver func(Conversation) bool, toolInputType string) Conversation {
	constructor := model.NewTextToolInput
	if toolInputType == "json" {
		constructor = model.NewJSONToolInput
	}
	c := ChatConversation{
		Base: Base{
			llm:              convoModel,
			isOver:           isOver,
			Messages:         initMessages,
			InputConstructor: constructor,
		},
	}

	for _, message := range c.Messages {
		message.Print()
	}

	return &c
}

func (c *ChatConversation) RunConversation() {
	for {
		response, err := c.generateModelResponse()
		if err != nil {
			slog.Error("LLM session failed", "err", err)
			break
		}

		input, err := c.InputConstructor(response)
		output := ""
		if err != nil {
			output = fmt.Sprintf("error: %s", err)
		} else {
			output = tools.RunTool(input)
		}

		c.Messages = append(c.Messages, model.NewChat("assistant", response))
		c.Messages = append(c.Messages, model.NewChat("user", output))

		c.Messages[len(c.Messages)-2].Print()
		c.Messages[len(c.Messages)-1].Print()

		if err == nil && c.isOver(c) {
			break
		}
	}
}

func (c *ChatConversation) generateModelResponse() (string, error) {
	raw, err := c.llm.CompleteChat(c.Messages)
	if err != nil {
		raw = ""
	}

	return raw, err
}
