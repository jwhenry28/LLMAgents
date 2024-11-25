package model

import "fmt"

const (
	JSON_TOOL_TYPE  = "json"
	TEXT_TOOL_TYPE  = "text"
	JSON_FORMAT_MSG = "Please respond in raw JSON format. Do not send any other text, including a markdown JSON code block."
	TEXT_FORMAT_MSG = "Respond using CLI format, e.g., <tool-name> <arg1> <arg2> ..."
)

type ToolInput interface {
	AsString() string
	GetName() string
	GetArgs() []string
	FormatUsage(string, []string) string
}

func NewToolInput(toolType, input string) (ToolInput, error) {
	switch toolType {
	case JSON_TOOL_TYPE:
		return NewJSONToolInput(input)
	case TEXT_TOOL_TYPE:
		return NewTextToolInput(input)
	default:
		return nil, fmt.Errorf("unknown tool-type %s", toolType)
	}
}