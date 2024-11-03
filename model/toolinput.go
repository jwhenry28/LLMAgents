package model

import "encoding/json"

type ToolInput struct {
	Name string   `json:"tool"`
	Args []string `json:"args"`
}

func ToolInputFromJSON(response string) (ToolInput, error) {
	var input ToolInput
	err := json.Unmarshal([]byte(response), &input)
	return input, err
}
