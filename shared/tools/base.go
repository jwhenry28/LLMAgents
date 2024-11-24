package tools

import "github.com/jwhenry28/LLMAgents/shared/model"

type Base struct {
	Input           model.ToolInput
	Name            string
	Args            []string
	BriefText       string
	ExplanationText string
}

func (task Base) Brief() string {
	return task.BriefText
}

func (task Base) Usage() string {
	return task.Input.FormatUsage(task.Args)
}

func (task Base) Explanation() string {
	return task.ExplanationText
}

func (task Base) Help() string {
	msg := task.Brief() + "\n" + task.Usage()
	if task.ExplanationText != "" {
		msg += "\n" + task.ExplanationText
	}

	return msg
}
