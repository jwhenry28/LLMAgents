package conversation

import (
	"fmt"
	"log/slog"

	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"

	"github.com/jwhenry28/LLMAgents/media-curator/tools"
)

type Conversation struct {
	llm      llm.LLM
	messages []model.Chat
}

func NewConversation(convoModel llm.LLM, initMessages []model.Chat) *Conversation {
	c := Conversation{
		llm:      convoModel,
		messages: initMessages,
	}

	for _, message := range c.messages {
		message.Print()
	}

	return &c
}

func (c *Conversation) RunConversation(seed string) {
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

		c.messages = append(c.messages, model.NewChat("assistant", response))
		c.messages = append(c.messages, model.NewChat("user", output))

		c.messages[len(c.messages)-2].Print()
		c.messages[len(c.messages)-1].Print()

		if c.isOver() {
			break
		}
	}
}

func (c *Conversation) generateModelResponse() (string, error) {
	raw, err := c.llm.CompleteChat(c.messages)
	if err != nil {
		raw = ""
	}

	return raw, err
}

func (c *Conversation) isOver() bool {
	modelResponse := c.messages[len(c.messages)-2]
	selectedTool, _ := model.ToolInputFromJSON(modelResponse.Content)

	toolOutput := c.messages[len(c.messages)-1]

	return selectedTool.Name == "decide" && (toolOutput.Content == "notified" || toolOutput.Content == "ignored")
}
