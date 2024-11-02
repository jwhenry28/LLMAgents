package curation

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"hackandpray.com/media-curator/llm"
	"hackandpray.com/media-curator/model"
	"hackandpray.com/media-curator/utils"
)

type Conversation struct {
	llm      llm.LLM
	ongoing  bool
	messages []model.Chat
}

func NewConversation(seed, description, initialHtml string) *Conversation {
	c := Conversation{}
	c.initialize(seed, description, initialHtml)
	return &c
}

func (c *Conversation) initialize(seed, description, initialHtml string) {
	c.llm = llm.NewMockLLM()
	c.messages = []model.Chat{}

	c.initializeMessages(seed, description, initialHtml)
}

func (c *Conversation) initializeMessages(seed, description, initialHtml string) {
	c.messages = []model.Chat{
		{
			Role:    "system",
			Content: fmt.Sprintf(utils.SYSTEM_PROMPT, description),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(utils.URL_PROMPT, seed, initialHtml),
		},
	}
}

func (c *Conversation) RunConversation(seed string) {
	c.ongoing = true
	for c.ongoing {
		c.generateModelResponse()
		c.runTool()
	}
}

func (c *Conversation) generateModelResponse() {
	raw, err := c.llm.CompleteChat(c.messages)
	var response model.Chat
	if err != nil {
		slog.Error("LLM session failed", "err", err)
		c.ongoing = false
		response = model.NewChat("", "")
	} else {
		response = model.NewChat("assistant", raw)
	}

	c.messages = append(c.messages, response)
}

func (c *Conversation) runTool() {
	input := c.messages[len(c.messages)-1].Content
	var tool model.Tool
	err := json.Unmarshal([]byte(input), &tool)

	output := "unknown tool: " + tool.Name
	if err != nil {
		output = fmt.Sprintf("LLM response is not conforming JSON: %s\nerror: %s\n", tool, err.Error())
	}

	if tool.Name == "decide" {
		args := tool.Args
		if len(args) < 2 {
			output = fmt.Sprintf("too few args to 'decide' call, expected 2: %v\n", args)
		}

		if args[0] == "NOTIFY" {
			slog.Info("Sending notification for url", "url", args[1])
			c.ongoing = false
			output = ""
		} else if args[0] == "IGNORE" {
			slog.Info("Ignoring url", "url", args[1])
			c.ongoing = false
			output = ""
		} else {
			output = "improperly formatted 'decide' call: " + args[0]
		}
	}

	c.messages = append(c.messages, model.NewChat("user", output))
}
