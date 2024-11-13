package conversation

import (
	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"
)

type CmdConversation struct {
	Base
}

func RunConversation(convoModel llm.LLM, initMessages []model.Chat, isOver func(Conversation) bool, toolInputType string) Conversation {
	constructor := model.NewTextToolInput
	if toolInputType == "json" {
		constructor = model.NewJSONToolInput
	}
	c := CmdConversation{
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

func (c *CmdConversation) RunConversation() {
	// Similar to ChatConversation, but do not send all messages to LLM at once

	// Instead, send one message at a time, and use the LLM to generate the next command

	/*
		You are a coder agent.
		You are given a goal, and you need to choose commands that gets you closer to the goal..
		You will probably not solve the goal with a single command.
		Consider that other agents will run after you, so it's okay if there is still other work to be done afterward.

	*/
}
