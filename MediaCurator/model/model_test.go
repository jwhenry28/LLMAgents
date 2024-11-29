package model

import (
	"os"
	"testing"
)

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

func TestNewProfileFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		want     *Profile
		wantErr  bool
	}{
		{
			name: "valid profile",
			yamlData: `name: Test User
email: test@example.com
interests: |
  Security research
  Machine learning`,
			want: &Profile{
				Name:      "Test User",
				Email:     "test@example.com",
				Interests: "Security research\nMachine learning",
			},
			wantErr: false,
		},
		{
			name:     "invalid yaml",
			yamlData: "invalid: [yaml: content",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with test YAML content
			tmpfile, err := os.CreateTemp("", "profile-*.yml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.yamlData)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Test the function
			got, err := NewProfileFromYAML(tmpfile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProfileFromYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.want.Name)
			}
			if got.Email != tt.want.Email {
				t.Errorf("Email = %q, want %q", got.Email, tt.want.Email)
			}
			if got.Interests != tt.want.Interests {
				t.Errorf("Interests = %q, want %q", got.Interests, tt.want.Interests)
			}
		})
	}
}
