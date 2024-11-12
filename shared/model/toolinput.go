package model

type ToolInput interface {
	AsString() string
	GetName() string
	GetArgs() []string
}
