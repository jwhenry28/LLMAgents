package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ToolInput interface {
	AsString() string
	GetName() string
	GetArgs() []string
}

type JSONToolInput struct {
	Name string   `json:"tool"`
	Args []string `json:"args"`
}

func NewJSONToolInput(response string) (ToolInput, error) {
	var input ToolInput
	err := json.Unmarshal([]byte(response), &input)
	return input, err
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

type TextToolInput struct {
	Name string
	Args []string
}

func (t TextToolInput) AsString() string {
	template := `COMMAND: %s

ARGS:
%s`
	return fmt.Sprintf(template, t.Name, strings.Join(t.Args, "\n"))
}

func (t TextToolInput) GetName() string {
	return t.Name
}

func (t TextToolInput) GetArgs() []string {
	return t.Args
}

// TOOD: handle multiline strings
func NewTextToolInput(response string) (ToolInput, error) {
	items := strings.Split(response, " ")
	t := &TextToolInput{Name: items[0], Args: items[1:]}
	return t, nil
}
