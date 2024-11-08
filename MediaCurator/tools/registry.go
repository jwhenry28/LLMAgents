package tools

import (
	"fmt"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Tool interface {
	Brief() string
	Help() string
	Match() bool
	Invoke() string
}

var Registry = map[string]func(model.ToolInput) Tool{
	"help":   NewHelp,
	"decide": NewDecide,
	"fetch":  NewFetch,
}

func RunTool(input model.ToolInput) string {
	constructor, ok := Registry[input.Name]
	if !ok {
		return fmt.Sprintf("unknown tool: %s. use 'help' tool to view supported tools", input.Name)
	}

	tool := constructor(input)
	if !tool.Match() {
		return fmt.Sprintf("improper usage of tool: %s\n%s", input.Name, tool.Help())
	}

	return tool.Invoke()
}
