package model

import (
	"reflect"
	"testing"
)

func TestNewChat(t *testing.T) {
	role := "user"
	content := "test message"

	chat := NewChat(role, content)

	if chat.Role != role {
		t.Errorf("Expected role %s but got %s", role, chat.Role)
	}
	if chat.Content != content {
		t.Errorf("Expected content %s but got %s", content, chat.Content)
	}
}

func TestNewAnchor(t *testing.T) {
	text := "Click here"
	href := "https://example.com"

	anchor := NewAnchor(text, href)

	if anchor.Text != text {
		t.Errorf("Expected text %s but got %s", text, anchor.Text)
	}
	if anchor.HRef != href {
		t.Errorf("Expected href %s but got %s", href, anchor.HRef)
	}
}

func TestToolInputAsString(t *testing.T) {
	tests := []struct {
		name     string
		input    ToolInput
		expected string
	}{
		{
			name: "tool with multiple args",
			input: ToolInput{
				Name: "decide",
				Args: []string{"NOTIFY", "https://example.com"},
			},
			expected: "decide NOTIFY https://example.com",
		},
		{
			name: "tool with single arg",
			input: ToolInput{
				Name: "download",
				Args: []string{"https://example.com"},
			},
			expected: "download https://example.com",
		},
		{
			name: "tool with no args",
			input: ToolInput{
				Name: "list",
				Args: []string{},
			},
			expected: "list",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.AsString()
			if got != test.expected {
				t.Errorf("Expected string %q but got %q", test.expected, got)
			}
		})
	}
}

func TestToolInputFromJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected ToolInput
		wantErr  bool
	}{
		{
			name: "valid json",
			json: `{"tool":"decide","args":["NOTIFY","https://example.com"]}`,
			expected: ToolInput{
				Name: "decide",
				Args: []string{"NOTIFY", "https://example.com"},
			},
			wantErr: false,
		},
		{
			name:     "invalid json",
			json:     `{"tool":}`,
			expected: ToolInput{},
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ToolInputFromJSON(test.json)

			if test.wantErr && err == nil {
				t.Errorf("ToolInputFromJSON() error = %v, wantErr %v", err, test.wantErr)
				return
			} else if !test.wantErr && !reflect.DeepEqual(got, test.expected) {
				t.Errorf("Expected %v but got %v", test.expected.AsString(), got.AsString())
			}
		})
	}
}
