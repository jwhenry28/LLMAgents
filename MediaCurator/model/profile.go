package model

import (
	"fmt"
	"net/mail"
	"os"

	"gopkg.in/yaml.v2"
)

type Profile struct {
	Name      string `yaml:"name"`
	Email     string `yaml:"email"`
	Interests string `yaml:"interests"`
}

func NewProfileFromYAML(filename string) (*Profile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile file: %v", err)
	}

	var profile Profile
	err = yaml.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse profile YAML: %v", err)
	}

	if profile.Name == "" || profile.Email == "" || profile.Interests == "" {
		return nil, fmt.Errorf("profile missing required fields")
	}

	_, err = mail.ParseAddress(profile.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email address: %v", err)
	}

	return &profile, nil
}
