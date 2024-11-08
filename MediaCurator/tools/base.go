package tools

import "hackandpray.com/media-curator/model"

type Base struct {
	Input     model.ToolInput
	BriefText string
	UsageText string
}

func (task Base) Brief() string {
	return task.BriefText
}

func (task Base) Help() string {
	return task.Brief() + "\n" + task.UsageText
}
