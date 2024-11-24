package llm

import (
	"os"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

const DEFAULT_RETRIES = 3

type LLM interface {
	CompleteChat([]model.Chat) (string, error)
	Type() string
}

func ConstructLLM(llmType string) LLM {
	switch llmType {
	case "human":
		return NewHuman()
	case "mock":
		return NewMockLLM()
	case "openai":
		return NewOpenAI(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"), 0)
	case "claude":
		fallthrough
	case "anthropic":
		return NewAnthropic(os.Getenv("ANTHROPIC_API_KEY"), os.Getenv("ANTHROPIC_MODEL"), 0)
	default:
		return nil
	}
}
