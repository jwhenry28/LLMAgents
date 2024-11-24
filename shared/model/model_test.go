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

func TestNewJSONToolInput(t *testing.T) {
	tests := []struct {
		name     string
		input    ToolInput
		expected string
	}{
		{
			name: "tool with multiple args",
			input: func() ToolInput {
				input, err := NewJSONToolInput(`{"tool":"decide","args":["NOTIFY","https://example.com"]}`)
				if err != nil {
					t.Errorf("Expected input to be non-nil: %v", err)
					return nil
				}
				return input
			}(),
			expected: "decide NOTIFY https://example.com",
		},
		{
			name: "tool with single arg",
			input: func() ToolInput {
				input, _ := NewJSONToolInput(`{"tool":"download","args":["https://example.com"]}`)
				return input
			}(),
			expected: "download https://example.com",
		},
		{
			name: "tool with no args",
			input: func() ToolInput {
				input, _ := NewJSONToolInput(`{"tool":"list","args":[]}`)
				return input
			}(),
			expected: "list",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.input == nil {
				t.Errorf("Expected input to be non-nil")
				return
			}
			got := test.input.AsString()
			if got != test.expected {
				t.Errorf("Expected string %q but got %q", test.expected, got)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected ToolInput
		wantErr  bool
	}{
		{
			name: "valid json",
			json: `{"tool":"decide","args":["NOTIFY","https://example.com"]}`,
			expected: func() ToolInput {
				input, _ := NewJSONToolInput(`{"tool":"decide","args":["NOTIFY","https://example.com"]}`)
				return input
			}(),
			wantErr: false,
		},
		{
			name: "invalid json",
			json: `{"tool":}`,
			expected: func() ToolInput {
				input, _ := NewJSONToolInput(`{"tool":}`)
				return input
			}(),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewJSONToolInput(test.json)

			if test.wantErr && err == nil {
				t.Errorf("FromString() error = %v, wantErr %v", err, test.wantErr)
				return
			} else if !test.wantErr && !reflect.DeepEqual(got, test.expected) {
				t.Errorf("Expected %v but got %v", test.expected.AsString(), got.AsString())
			}
		})
	}
}

func TestNewTextToolInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ToolInput
		wantErr  bool
	}{
		{
			name:  "simple command",
			input: "decide NOTIFY https://example.com",
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"NOTIFY", "https://example.com"},
			},
			wantErr: false,
		},
		{
			name:  "quoted args",
			input: `decide "hello world" 'test quote'`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"hello world", "test quote"},
			},
			wantErr: false,
		},
		{
			name:  "tricky quoted args",
			input: `decide "hello world' 'test quote"`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"hello world' 'test quote"},
			},
			wantErr: false,
		},
		{
			name:  "nested quoted args",
			input: `decide "hello world" 'test quote "nested quote"'`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"hello world", "test quote \"nested quote\""},
			},
			wantErr: false,
		},
		{
			name:  "multiline input",
			input: "decide This\nis a\nmultiline argument",
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"This", "is a\nmultiline argument"},
			},
			wantErr: false,
		},
		{
			name:  "multiline input2",
			input: "decide\nThis\nis another\nmultiline argument",
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"This\nis another\nmultiline argument"},
			},
			wantErr: false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "nested quotes with escape",
			input: `decide "hello \"world\" and 'test quote'"`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{`hello "world" and 'test quote'`},
			},
			wantErr: false,
		},
		{
			name:  "multiple spaces",
			input: `decide     "hello     world" 'test quote'`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{`hello     world`, `test quote`},
			},
			wantErr: false,
		},
		{
			name:  "unbalanced quotes",
			input: `decide "hello world'`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{`hello world'`}, // This should still be treated as a single argument
			},
			wantErr: false,
		},
		{
			name:     "empty input with spaces",
			input:    "     ",
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "only quotes",
			input: `'"' "'"`,
			expected: &TextToolInput{
				Name: `"`,
				Args: []string{`'`},
			},
			wantErr: false,
		},
		{
			name:  "mixed quotes",
			input: `decide "hello 'world'"`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{`hello 'world'`},
			},
			wantErr: false,
		},
		{
			name:  "complex nested quotes",
			input: `decide "this is a 'complex' example with \"nested quotes\""`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{`this is a 'complex' example with "nested quotes"`},
			},
			wantErr: false,
		},
		{
			name: "newline with quotes",
			input: `decide "this is a 'complex'
"example with
\"nested quotes\"
and a newline"`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"this is a 'complex'", `example with
\"nested quotes\"
and a newline`},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewTextToolInput(test.input)

			if test.wantErr && err == nil {
				t.Errorf("NewTextToolInput() error = %v, wantErr %v", err, test.wantErr)
				return
			} else if !test.wantErr && !reflect.DeepEqual(got, test.expected) {
				t.Errorf("Expected\n%v\n\nbut got\n\n%v", test.expected.AsString(), got.AsString())
			}
		})
	}
}

func TestNewToolInput(t *testing.T) {
	tests := []struct {
		name     string
		toolType string
		input    string
		expected ToolInput
		wantErr  bool
	}{
		{
			name:     "valid json tool",
			toolType: JSON_TOOL_TYPE,
			input:    `{"tool": "decide", "args": ["hello", "world"]}`,
			expected: &JSONToolInput{
				Name: "decide",
				Args: []string{"hello", "world"},
			},
			wantErr: false,
		},
		{
			name:     "valid text tool",
			toolType: TEXT_TOOL_TYPE,
			input:    `decide hello world`,
			expected: &TextToolInput{
				Name: "decide",
				Args: []string{"hello", "world"},
			},
			wantErr: false,
		},
		{
			name:     "invalid tool type",
			toolType: "invalid",
			input:    "anything",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid json input",
			toolType: JSON_TOOL_TYPE,
			input:    `{"tool": "decide"`, // malformed JSON
			expected: nil,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewToolInput(test.toolType, test.input)

			if test.wantErr {
				if err == nil {
					t.Errorf("NewToolInput() error = nil, wantErr %v", test.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewToolInput() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("NewToolInput() = %v, want %v", got, test.expected)
			}
		})
	}
}
