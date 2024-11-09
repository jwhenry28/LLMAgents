package tools

import (
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Finish struct {
	AllowedArgs []string
	tools.Base
}

func NewFinish(input model.ToolInput) tools.Tool {
	brief := "finish: let the user know you are finished. note, running this tool will end the conversation."
	usage := `usage: { "tool": "finish", "args": []}`
	return Finish{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Finish) Match() bool {
	return true
}

func (task Finish) Invoke() string {
	return ""
}
