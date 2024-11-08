package llm

import (
	"os"

	"hackandpray.com/media-curator/model"
)

type LLM interface {
	CompleteChat([]model.Chat) (string, error)
}

func ConstructLLM(llmType string) LLM {
	switch llmType {
	case "human":
		return NewHuman()
	case "mock":
		return NewMockLLM()
	case "openai":
		return NewOpenAI(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"), 0)
	default:
		return nil
	}
}