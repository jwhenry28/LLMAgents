package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
)

type Write struct {
	AllowedArgs []string
	tools.Base
}

func NewWrite(input model.ToolInput) tools.Tool {
	brief := "write: writes text to a file."
	usage := `usage: { "tool": "write", "args": [ <filename>, <text> ]}
args:
- filename: The name of the file to write to
- text: The text to write to the file`
	return Write{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Write) Match() bool {
	return len(task.Input.Args) >= 2
}

func (task Write) Invoke() string {
	args := task.Input.Args
	filename := args[0]
	text := args[1]

	if strings.HasPrefix(filename, "/") {
		return "error: must use absolute paths"
	}

	if strings.Contains(filename, "..") {
		return "error: cannot use traversal character (..)"
	}

	err := os.WriteFile(utils.SANDBOX_DIR + "/" + filename, []byte(text), 0644)
	if err != nil {
		return fmt.Sprintf("error writing to file: %s", err)
	}

	return "success"
}
