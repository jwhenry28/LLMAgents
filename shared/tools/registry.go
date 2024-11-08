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

var Registry = make(map[string]func(model.ToolInput) Tool)

func RegisterTool(name string, constructor func(model.ToolInput) Tool) error {
	if _, ok := Registry[name]; ok {
		return fmt.Errorf("tool already registered: %s", name)
	}

	Registry[name] = constructor
	return nil
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
