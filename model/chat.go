package model

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
