package tools

import (
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Report struct {
	tools.Base
}

func NewReport(input model.ToolInput) tools.Tool {
	brief := "report: share your final report and end the conversation. only use this tool when you have completed the task and want to share the final output."
	usage := `usage: report <report_text>
args:
- report_text: The text of the report to share`
	return Report{
		Base: tools.Base{Input: input, BriefText: brief, UsageText: usage},
	}
}

func (task Report) Match() bool {
	args := task.Input.GetArgs()
	return len(args) > 0
}

func (task Report) Invoke() string {
	args := task.Input.GetArgs()
	return args[0]
}
