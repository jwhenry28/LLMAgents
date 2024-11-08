package tools

import (
	"golang.org/x/exp/slices"
	"hackandpray.com/media-curator/model"
)

type Decide struct {
	AllowedArgs []string
	Base
}

func NewDecide(input model.ToolInput) Tool {
	brief := "decide: issues a final decision on the specified URL."
	usage := `usage: { "tool": "decide", "args": [ <decision>, <url>, <justification> ]}
args:
- url: The URL you are making a decision about
- decision: Your decision. Must be one of the following:
	- IGNORE: Choose this option if you do not think your client will be interested in reading this URL today.
	- NOTIFY: Choose this option if you would like to forward this URL to your client
- justification: (optional) A short explanation for your decision`
	return Decide{
		AllowedArgs: []string{"NOTIFY", "IGNORE"},
		Base:        Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Decide) Match() bool {
	return len(task.Input.Args) >= 2 && slices.Contains(task.AllowedArgs, task.Input.Args[0])
}

func (task Decide) Invoke() string {
	args := task.Input.Args
	if args[0] == "NOTIFY" {
		return "notified"
	} else if args[0] == "IGNORE" {
		return "ignored"
	}

	return "unknown decision"
}
