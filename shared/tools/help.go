package tools

import (
	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Help struct {
	Base
}

func NewHelp(input model.ToolInput) Tool {
	name := "help"
	args := []string{"tool-name"}
	brief := "help: returns information about supported tools. If no arguments are supplied, returns a list of all tool names. If a tool name is supplied as an argument, retrieved specific information about that tool."
	explanation := `args: 
- tool-name: optional argument. if included, this specifies one tool to learn more about`
	return Help{
		Base: Base{Input: input, Name: name, Args: args, BriefText: brief, ExplanationText: explanation},
	}
}

func (task Help) Match() bool {
	return true
}

func (task Help) Invoke() string {
	args := task.Input.GetArgs()
	output := ""
	if len(args) == 0 {
		output = GetToolList()
	} else {
		output = GetToolHelp(args[0])
	}

	return output
}

func GetToolList() string {
	output := ""
	for _, constructor := range Registry {
		output += " - " + constructor(model.TextToolInput{}).Brief() + "\n"
	}

	return output
}

func GetToolHelp(toolName string) string {
	constructor, ok := Registry[toolName]
	output := ""
	if !ok {
		output = "unknown tool: %s. supported tools:\n"
		output += GetToolList()
	} else {
		output = constructor(model.TextToolInput{}).Help()
	}

	return output
}
