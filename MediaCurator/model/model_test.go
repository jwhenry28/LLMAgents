package model

import (
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
