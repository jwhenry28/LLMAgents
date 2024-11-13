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
	usage := `usage: write <filename>
<text>
args:
- filename: The name of the file to write to
- text: The text to write to the file. Note, use actual newlines, not \\n.`
	return Write{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Write) Match() bool {
	args := task.Input.GetArgs()
	return len(args) >= 2
}

func (task Write) Invoke() string {
	args := task.Input.GetArgs()
	filename := args[0]
	text := args[1]

	if strings.HasPrefix(filename, "/") {
		return "error: must use absolute paths"
	}

	if strings.Contains(filename, "..") {
		return "error: cannot use traversal character (..)"
	}
	// Create all necessary directories
	dir := utils.SANDBOX_DIR + "/" + filename
	if lastSlash := strings.LastIndex(dir, "/"); lastSlash != -1 {
		dirPath := dir[:lastSlash]
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Sprintf("error creating directories: %s", err)
		}
	}

	err := os.WriteFile(utils.SANDBOX_DIR+"/"+filename, []byte(text), 0644)
	if err != nil {
		return fmt.Sprintf("error writing to file: %s", err)
	}

	return "success"
}
