package llm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Human struct {
}

func NewHuman() *Human {
	return &Human{}
}

func (llm *Human) CompleteChat(messages []model.Chat) (string, error) {
	tool := ""
	args := []string{}

	fmt.Println(messages[len(messages)-2])
	fmt.Println(messages[len(messages)-1])

	fmt.Print("Enter tool and args (space separated): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.Fields(scanner.Text())
		if len(input) > 0 {
			tool = input[0]
			if len(input) > 1 {
				args = input[1:]
			}
		}
	}

	response := `{ "tool": "` + tool + `", "args": [ `
	for i, arg := range args {
		response += `"` + arg + `"`
		if i < len(args)-1 {
			response += `, `
		}
	}
	response += `] }`

	return response, nil
}
