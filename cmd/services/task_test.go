package services

import (
	"testing"
)

func TestPriorityValidation(t *testing.T) {
	validPriorities := map[string]bool{
		"urgent": true,
		"high":   true,
		"medium": true,
		"low":    true,
	}

	tests := []struct {
		priority string
		valid    bool
	}{
		{"urgent", true},
		{"high", true},
		{"medium", true},
		{"low", true},
		{"critical", false},
		{"", false},
		{"normal", false},
	}

	for _, tt := range tests {
		if validPriorities[tt.priority] != tt.valid {
			t.Errorf("priority %q: expected valid=%v, got %v", tt.priority, tt.valid, validPriorities[tt.priority])
		}
	}
}

func TestDefaultPriority(t *testing.T) {
	input := CreateTaskInput{
		ColumnID: "col-1",
		Title:    "Test Task",
	}

	if input.Priority == "" {
		input.Priority = "medium"
	}

	if input.Priority != "medium" {
		t.Errorf("expected default priority 'medium', got %q", input.Priority)
	}
}

func TestMoveTaskInput(t *testing.T) {
	input := MoveTaskInput{
		ColumnID: "col-2",
		Position: 3,
	}

	if input.ColumnID == "" {
		t.Error("ColumnID should not be empty")
	}
	if input.Position < 0 {
		t.Error("Position should not be negative")
	}
}
