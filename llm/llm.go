package llm

import (
	"hackandpray.com/media-curator/model"
)

type LLM interface {
	CompleteChat([]model.Chat) (string, error)
}
