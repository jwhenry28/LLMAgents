package tools

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
)

type Run struct {
	AllowedArgs []string
	tools.Base
}

func NewRun(input model.ToolInput) tools.Tool {
	brief := "run: Runs a Golang program."
	usage := `usage: run <filename>
args:
- filename: The name of the file to run`
	return Run{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Run) Match() bool {
	args := task.Input.GetArgs()
	return len(args) >= 1
}

func (task Run) Invoke() string {
	args := task.Input.GetArgs()
	filename := args[0]

	if strings.HasPrefix(filename, "/") {
		return "error: must use absolute paths"
	}

	if strings.Contains(filename, "..") {
		return "error: cannot use traversal character (..)"
	}

	cmd := exec.Command("go", "run", utils.SANDBOX_DIR + "/" + filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error running program: %s\n%s", err, string(output))
	}
	return string(output)
}
