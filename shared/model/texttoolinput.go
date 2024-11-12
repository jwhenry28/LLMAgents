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
	if strings.Contains(response, "\n") { // Check for multiple lines
		items := strings.Fields(response) // Split by whitespace
		if len(items) > 0 {
			name := items[0]
			args := strings.TrimSpace(response[len(name):]) // Everything after the first space
			return &TextToolInput{Name: name, Args: []string{args}}, nil
		}
		return nil, fmt.Errorf("invalid response format")
	}

	// Command line parsing with quotes
	parsedArgs := parseCommandLine(response) // Use a helper function to parse command line
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

// Helper function to parse command line input
func parseCommandLine(input string) []string {
	var args []string
	quoted := false
	quoteChar := rune(0)
	args = strings.FieldsFunc(input, func(r rune) bool {
		ok := !quoted && r == ' '
		if r == '"' || r == '\'' {
			if quoted && r == quoteChar {
				ok = quoted && r == quoteChar
				quoted = false
				quoteChar = rune(0)
			} else if !quoted {
				quoted = true
				quoteChar = r
				ok = true
			}
		}
		return ok
	})

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
