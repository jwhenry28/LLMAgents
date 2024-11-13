package tools

import (
	"fmt"
	"os/exec"

	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type GoInit struct {
	tools.Base
}

func NewGoInit(input model.ToolInput) tools.Tool {
	brief := "goinit: initializes a new Go module"
	usage := `usage: goinit <module name>
args:
	- module name: The name of the module to initialize`
	return GoInit{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task GoInit) Match() bool {
	args := task.Input.GetArgs()
	return len(args) == 1
}

func (task GoInit) Invoke() string {
	args := task.Input.GetArgs()
	cmd := exec.Command("go", "mod", "init", args[0])
	cmd.Dir = utils.SANDBOX_DIR
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error running program: %s\n%s", err, string(output))
	}
	return string(output)
}
