package tools

import (
	"golang.org/x/exp/slices"

	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Decide struct {
	AllowedArgs []string
	tools.Base
}

func NewDecide(input model.ToolInput) tools.Tool {
	brief := "decide: issues a final decision on the specified URL."
	usage := `usage: { "tool": "decide", "args": [ <decision>, <url>, <justification> ]}
args:
- url: The URL you are making a decision about
- decision: Your decision. Must be one of the following:
	- IGNORE: Choose this option if you do not think your client will be interested in reading this URL today.
	- NOTIFY: Choose this option if you would like to forward this URL to your client
- justification: A short explanation for your decision`
	return Decide{
		AllowedArgs: []string{"NOTIFY", "IGNORE"},
		Base:        tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Decide) Match() bool {
	args := task.Input.GetArgs()
	return len(args) >= 3 && slices.Contains(task.AllowedArgs, args[0])
}

func (task Decide) Invoke() string {
	args := task.Input.GetArgs()
	if args[0] == "NOTIFY" {
		return "notified"
	} else if args[0] == "IGNORE" {
		return "ignored"
	}

	return "unknown decision"
}
