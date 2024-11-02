package model

type Tool struct {
	Name string   `json:"tool"`
	Args []string `json:"args"`
}
