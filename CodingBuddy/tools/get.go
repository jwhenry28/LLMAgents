package tools

import (
	"fmt"
	"os/exec"

	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type GoGet struct {
	tools.Base
}

func NewGoGet(input model.ToolInput) tools.Tool {
	brief := "goget: downloads a Go module (akin to 'go get')"
	usage := `usage: { "tool": "goget", "args": [ <module name> ]}
args:
	- module name: The name of the module to download`
	return GoGet{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task GoGet) Match() bool {
	return len(task.Input.Args) >= 1
}

func (task GoGet) Invoke() string {
	args := task.Input.Args
	moduleName := args[0]

	cmd := exec.Command("go", "get", moduleName)
	cmd.Dir = utils.SANDBOX_DIR
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error running program: %s\n%s", err, string(output))
	}
	return string(output)
}
