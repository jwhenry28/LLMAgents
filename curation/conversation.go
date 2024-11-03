package curation

import (
	"fmt"
	"log/slog"

	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
	"hackandpray.com/media-curator/tools"
)

type Conversation struct {
	llm      llm.LLM
	ongoing  bool
	messages []model.Chat
}

func NewConversation(convoModel llm.LLM, initMessages []model.Chat) *Conversation {
	c := Conversation{
		llm:      convoModel,
		ongoing:  false,
		messages: initMessages,
	}

	for _, message := range initMessages {
		fmt.Println(message)
	}

	return &c
}

func (c *Conversation) RunConversation(seed string) {
	c.ongoing = true
	for c.ongoing {
		response, err := c.generateModelResponse()
		if err != nil {
			break
		}

		responseChat := model.NewChat("assistant", response)
		fmt.Println(responseChat)

		input, err := model.ToolInputFromJSON(response)
		c.ongoing = !c.isOver(input, err)

		output := tools.RunTool(input)

		outputChat := model.NewChat("user", output)
		fmt.Println(outputChat)

		c.messages = append(c.messages, responseChat)
		c.messages = append(c.messages, outputChat)
	}
}

func (c *Conversation) generateModelResponse() (string, error) {
	raw, err := c.llm.CompleteChat(c.messages)
	if err != nil {
		slog.Error("LLM session failed", "err", err)
		raw = ""
	}

	return raw, err
}

func (c *Conversation) isOver(selectedTool model.ToolInput, previousError error) bool {
	return selectedTool.Name == "decide" || previousError != nil
}
