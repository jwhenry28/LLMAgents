package conversation

import (
	"fmt"
	"log/slog"

	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"

	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Conversation struct {
	llm    llm.LLM
	isOver func(*Conversation) bool

	Messages []model.Chat
}

func NewConversation(convoModel llm.LLM, initMessages []model.Chat, isOver func(*Conversation) bool) *Conversation {
	c := Conversation{
		llm:      convoModel,
		isOver:   isOver,
		Messages: initMessages,
	}

	for _, message := range c.Messages {
		message.Print()
	}

	return &c
}

func (c *Conversation) RunConversation() {
	for {
		response, err := c.generateModelResponse()
		if err != nil {
			slog.Error("LLM session failed", "err", err)
			break
		}

		input, err := model.ToolInputFromJSON(response)
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

		if c.isOver(c) {
			break
		}
	}
}

func (c *Conversation) generateModelResponse() (string, error) {
	raw, err := c.llm.CompleteChat(c.Messages)
	if err != nil {
		raw = ""
	}

	return raw, err
}
