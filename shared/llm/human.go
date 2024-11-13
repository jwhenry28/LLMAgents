package llm

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Human struct {
}

func NewHuman() *Human {
	return &Human{}
}

func (llm *Human) CompleteChat(messages []model.Chat) (string, error) {
	fmt.Print("\nEnter tool and args (space separated):\n")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Look for three consecutive newlines
		for i := 0; i < len(data)-2; i++ {
			if data[i] == '\n' && data[i+1] == '\n' && data[i+2] == '\n' {
				return i + 3, data[0:i], nil
			}
		}

		// If we're at EOF, return all remaining data
		if atEOF {
			return len(data), data, nil
		}

		// Request more data
		return 0, nil, nil
	})

	input := ""
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
		input += scanner.Text()
	}

	return input, nil
}
