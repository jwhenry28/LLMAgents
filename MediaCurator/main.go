package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jwhenry28/LLMAgents/media-curator/curation"
	"github.com/jwhenry28/LLMAgents/media-curator/model"
	"github.com/jwhenry28/LLMAgents/shared/llm"
)

const PROFILES_DIR = "data/profiles"

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
	sendEmailFlag := flag.Bool("send-email", false, "Send results to email address (optional)")
	profileFlag := flag.String("profile", "", "Specific profile to run (optional)")
	flag.Parse()

	llm := llm.ConstructLLM(*llmTypeFlag)
	if llm == nil {
		slog.Error("failed to create llm", "type", *llmTypeFlag)
		return
	}

	profiles := make(map[string]*model.Profile)
	if *profileFlag != "" {
		profile, err := loadSingleProfile(*profileFlag)
		if err != nil {
			slog.Error("error loading profile", "error", err)
			return
		}
		profiles[profile.Email] = profile
	} else {
		profiles, err = loadProfiles()
		if err != nil {	
			slog.Error("error loading profiles", "error", err)
			return
		}
	}

	for _, profile := range profiles {
		curator := curation.NewCurator(llm, profile)
		curator.SetSendEmail(*sendEmailFlag)
		curator.Curate()
	}
}

func loadSingleProfile(name string) (*model.Profile, error) {
	profilePath := fmt.Sprintf("%s/%s.yml", PROFILES_DIR, name)
	return model.NewProfileFromYAML(profilePath)
}

func loadProfiles() (map[string]*model.Profile, error) {
	files, err := os.ReadDir(PROFILES_DIR)
	if err != nil {
		slog.Error("error reading profiles directory", "error", err)
		return nil, err
	}

	profiles := make(map[string]*model.Profile)
	for _, file := range files {
		profilePath := fmt.Sprintf("%s/%s", PROFILES_DIR, file.Name())
		profile, err := model.NewProfileFromYAML(profilePath)
		if err != nil {
			slog.Warn("error loading profile", "file", file.Name(), "error", err)
			continue
		}
		profiles[profile.Email] = profile
	}

	return profiles, nil
}
