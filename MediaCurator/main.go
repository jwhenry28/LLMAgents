package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/jwhenry28/LLMAgents/media-curator/curation"
	"github.com/jwhenry28/LLMAgents/shared/llm"
)

func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return scanner.Err()
}

func main() {
	err := loadEnv()
	if err != nil {
		slog.Warn("Error loading environment variables", "error", err)
	}

	llmTypeFlag := flag.String("llm", "", "Type of LLM to use (openai, human, mock)")
	emailRecipient := flag.String("email", "", "Email address to send results to (optional)")
	flag.Parse()

	llm := llm.ConstructLLM(*llmTypeFlag)
	if llm == nil {
		slog.Error("failed to create llm", "type", *llmTypeFlag)
		return
	}
	curator := curation.NewCurator(llm, *emailRecipient)
	curator.Curate()
}
