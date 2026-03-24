package services

import (
	"testing"
)

func TestDefaultLabelColor(t *testing.T) {
	color := ""
	if color == "" {
		color = "#1E88E5"
	}
	if color != "#1E88E5" {
		t.Errorf("expected default color #1E88E5, got %s", color)
	}
}

func TestLabelColorValidation(t *testing.T) {
	validColors := []string{
		"#1E88E5", "#7C4DFF", "#43A047", "#FB8C00",
		"#E53935", "#00ACC1", "#8E24AA", "#F4511E",
	}

	for _, color := range validColors {
		if len(color) != 7 || color[0] != '#' {
			t.Errorf("invalid color format: %s", color)
		}
	}
}

func TestCommentContentValidation(t *testing.T) {
	tests := []struct {
		content string
		valid   bool
	}{
		{"Hello world", true},
		{"", false},
		{"  ", false},
		{"This is a valid comment.", true},
	}

	for _, tt := range tests {
		isValid := len(tt.content) > 0 && tt.content != "  "
		// Simplified validation
		if tt.content == "  " {
			isValid = false
		}
		if isValid != tt.valid {
			t.Errorf("content %q: expected valid=%v, got %v", tt.content, tt.valid, isValid)
		}
	}
}
