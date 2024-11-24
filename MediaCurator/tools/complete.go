package tools

import (
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Complete struct {
	AllowedArgs []string
	tools.Base
}

func NewComplete(input model.ToolInput) tools.Tool {
	name := "complete"
	args := []string{}
	brief := "complete: ends the conversation"
	explanation := `args:
- none
`
	return Complete{
		AllowedArgs: []string{"NOTIFY", "IGNORE"},
		Base:        tools.Base{Input: input, Name: name, Args: args, BriefText: brief, ExplanationText: explanation},
	}
}

func (task Complete) Match() bool {
	return len(task.Input.GetArgs()) == 0
}

func (task Complete) Invoke() string {
	return "completed"
}
