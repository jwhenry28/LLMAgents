package model

import (
	"fmt"
	"strings"
)

type TextToolInput struct {
	Name string
	Args []string
}

func NewTextToolInput(response string) (ToolInput, error) {
	if strings.Contains(response, "\n") {
		return handleMultiline(response)
	}
	return handleSingleLine(response)
}

func handleMultiline(response string) (ToolInput, error) {
	items := strings.Fields(response)
	if len(items) > 0 {
		name := items[0]
		args := strings.TrimSpace(response[len(name):])
		trimQuotes := (args[0] == '"' || args[0] == '\'') && args[0] == args[len(args)-1]
		
		if trimQuotes {
			args = args[1 : len(args)-1]
		}
		return &TextToolInput{Name: name, Args: []string{args}}, nil
	}
	return nil, fmt.Errorf("invalid response format")
}

func handleSingleLine(response string) (ToolInput, error) {
	parsedArgs := parseCommandLine(response)
	if len(parsedArgs) == 0 {
		return nil, fmt.Errorf("invalid response format")
	}
	name := parsedArgs[0]
	args := []string{}
	if len(parsedArgs) > 1 {
		args = parsedArgs[1:]
	}

	return &TextToolInput{Name: name, Args: args}, nil
}

func parseCommandLine(input string) []string {
	quoted := false
	quoteChar := rune(0)
	previousChar := rune(0)

	startQuoting := func(r rune) bool {
		return (r == '"' || r == '\'') && !quoted && previousChar != '\\'
	}

	endQuoting := func(r rune) bool {
		return (r == '"' || r == '\'') && quoted && r == quoteChar && previousChar != '\\'
	}

	shouldParse := func(r rune) bool {
		parse := !quoted && r == ' '

		if startQuoting(r) {
			parse = true
			quoted = true
			quoteChar = r
		} else if endQuoting(r) {
			parse = true
			quoted = false
			quoteChar = rune(0)
		}

		previousChar = r
		return parse
	}

	args := strings.FieldsFunc(input, shouldParse)
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
		args[i] = strings.ReplaceAll(args[i], "\\\"", "\"")
		args[i] = strings.ReplaceAll(args[i], "\\'", "'")
	}
	return args
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
