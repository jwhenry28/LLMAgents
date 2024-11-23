package llm

import (
	"bufio"
	"bytes"
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
	fmt.Println("Enter your multi-line input (press Ctrl+D on Unix or Ctrl+Z on Windows when done):")

	reader := bufio.NewReader(os.Stdin)
	var buffer bytes.Buffer

	defer fmt.Println("")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Return accumulated input without the error if we hit EOF
				return buffer.String(), nil
			}
			// Return any other error
			return buffer.String(), err
		}
		buffer.WriteString(line)
	}
}
