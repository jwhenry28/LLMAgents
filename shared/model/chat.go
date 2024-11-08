package model

import (
	"fmt"
	"strings"
)

type Chat struct {
	Role    string `json:"role"     binding:"required"`
	Content string `json:"content"  binding:"required"`
}

func NewChat(role, content string) Chat {
	return Chat{
		Role:    role,
		Content: content,
	}
}

func (c Chat) Print() {
	fmt.Println(strings.ToUpper(c.Role))
	fmt.Println(c.Content)
	fmt.Println("------------------------------------------------------------------")

}
