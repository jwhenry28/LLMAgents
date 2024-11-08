package tools

import (
	"hackandpray.com/media-curator/model"
)

type Help struct {
	Base
}

func NewHelp(input model.ToolInput) Tool {
	brief := "help: returns information about supported tools. If no arguments are supplied, returns a list of all tool names. If a tool name is supplied as an argument, retrieved specific information about that tool."
	usage := `usage: { "tool": "help", "args": [ <tool-name> ]}
args: 
- tool-name: optional argument. if included, this specifies one tool to learn more about`
	return Help{
		Base: Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Help) Match() bool {
	return true
}

func (task Help) Invoke() string {
	args := task.Input.Args
	output := ""
	if len(args) == 0 {
		output = getToolList()
	} else {
		output = getToolHelp(args[0])
	}

	return output
}

func getToolList() string {
	output := ""
	for _, constructor := range Registry {
		output += " - " + constructor(model.ToolInput{}).Brief() + "\n"
	}

	return output
}

func getToolHelp(toolName string) string {
	constructor, ok := Registry[toolName]
	output := ""
	if !ok {
		output = "unknown tool: %s. supported tools:\n"
		output += getToolList()
	} else {
		output = constructor(model.ToolInput{}).Help()
	}

	return output
}
