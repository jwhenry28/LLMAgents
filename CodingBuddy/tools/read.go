package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Read struct {
	AllowedArgs []string
	tools.Base
}

func NewRead(input model.ToolInput) tools.Tool {
	brief := "read: reads text from a file."
	usage := `usage: read <filename>
args:
- filename: The name of the file to read`
	return Read{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Read) Match() bool {
	args := task.Input.GetArgs()
	return len(args) >= 1
}

func (task Read) Invoke() string {
	args := task.Input.GetArgs()
	filename := args[0]

	if strings.HasPrefix(filename, "/") {
		return "error: must use absolute paths"
	}

	if strings.Contains(filename, "..") {
		return "error: cannot use traversal character (..)"
	}

	content, err := os.ReadFile(utils.SANDBOX_DIR + "/" + filename)
	if err != nil {
		return fmt.Sprintf("error reading file: %s", err)
	}

	return string(content)
}
