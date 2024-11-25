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
	response = strings.TrimSpace(response)
	if strings.Contains(response, "\n") {
		return handleMultiline(response)
	}
	return handleSingleLine(response)
}

func handleMultiline(response string) (ToolInput, error) {
	lines := strings.SplitN(response, "\n", 2)
	if len(lines) == 0 {
		return &TextToolInput{}, fmt.Errorf("invalid response format")
	}

	firstLine := parseCommandLine(lines[0])
	if len(firstLine) == 0 {
		return &TextToolInput{}, fmt.Errorf("invalid response format") 
	}

	name := firstLine[0]
	args := []string{}
	if len(firstLine) > 1 {
		args = firstLine[1:]
	}

	if len(lines) > 1 {
		firstChar := lines[1][0]
		lastChar := lines[1][len(lines[1])-1]
		if (firstChar == '"' || firstChar == '\'') && firstChar == lastChar {
			lines[1] = lines[1][1 : len(lines[1])-1]
		}
		args = append(args, lines[1])
	}

	return &TextToolInput{Name: name, Args: args}, nil
}

func handleSingleLine(response string) (ToolInput, error) {
	parsedArgs := parseCommandLine(response)
	if len(parsedArgs) == 0 {
		return &TextToolInput{}, fmt.Errorf("invalid response format")
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

func (t TextToolInput) FormatUsage(name string, argNames []string) string {
	template := `usage: %s %s`
	argString := ""
	for i, arg := range argNames {
		argString += fmt.Sprintf("<%s>", arg)
		if i < len(argNames)-1 {
			argString += " "
		}
	}

	return fmt.Sprintf(template, name, argString)
}