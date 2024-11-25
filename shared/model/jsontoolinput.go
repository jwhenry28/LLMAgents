package model

import (
	"encoding/json"
	"fmt"
)

type JSONToolInput struct {
	Name string   `json:"tool"`
	Args []string `json:"args"`
}


func NewJSONToolInput(response string) (ToolInput, error) {
	var input JSONToolInput
	err := json.Unmarshal([]byte(response), &input)
	return &input, err
}

func (t JSONToolInput) AsString() string {
	output := t.Name
	for _, arg := range t.Args {
		output += " " + arg
	}
	return output
}

func (t JSONToolInput) GetName() string {
	return t.Name
}

func (t JSONToolInput) GetArgs() []string {
	return t.Args
}

func (t JSONToolInput) FormatUsage(name string, argNames []string) string {
	template := `usage: { "tool": "%s", "args": [ %s ]}`
	argString := ""
	for i, arg := range argNames {
		argString += fmt.Sprintf("<%s>", arg)
		if i < len(argNames)-1 {
			argString += ", "
		}
	}

	return fmt.Sprintf(template, name, argString)
}
