package tools

import (
	"log/slog"

	"golang.org/x/exp/slices"
	"hackandpray.com/media-curator/model"
)

type Decide struct {
	AllowedArgs []string
	Base
}

func NewDecide(input model.ToolInput) Tool {
	return Decide{
		AllowedArgs: []string{"NOTIFY", "IGNORE"},
		Base: Base{Input: input},
	}
}

func (task Decide) Help() string {
	help := `
decide: issues a final decision on the specified URL.
usage: { "tool": "decide", "args": [ <decision>, <url> ]}
args:
- url: The URL you are making a decision about
- decision: Your decision. Must be one of the following:
	- IGNORE: Choose this option if you do not think your client will be interested in reading this URL today.
	- NOTIFY: Choose this option if you would like to forward this URL to your client`
	return help
}

func (task Decide) Match() bool {
	return len(task.Input.Args) >= 2 && slices.Contains(task.AllowedArgs, task.Input.Args[0])
}

func (task Decide) Invoke() string {
	args := task.Input.Args
	if args[0] == "NOTIFY" {
		slog.Info("Sending notification for url", "url", args[1])
	} else if args[0] == "IGNORE" {
		slog.Info("Ignoring url", "url", args[1])
	}

	return ""
}
