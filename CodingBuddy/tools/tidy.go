package tools

import (
	"fmt"
	"os/exec"

	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type GoTidy struct {
	tools.Base
}

func NewGoTidy(input model.ToolInput) tools.Tool {
	brief := "gotidy: tidies a Go module (akin to 'go mod tidy')"
	usage := `usage: gotidy
args:
	- none`
	return GoTidy{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task GoTidy) Match() bool {
	args := task.Input.GetArgs()
	return len(args) == 0
}

func (task GoTidy) Invoke() string {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = utils.SANDBOX_DIR
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error running program: %s\n%s", err, string(output))
	}
	return string(output)
}
