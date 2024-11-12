package model

import "encoding/json"

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
