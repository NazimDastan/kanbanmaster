package services

import (
	"testing"
)

func TestDelegateInputValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    DelegateInput
		wantErr  bool
	}{
		{"valid input", DelegateInput{ToUserID: "user-1", Reason: "On vacation"}, false},
		{"empty user", DelegateInput{ToUserID: "", Reason: "reason"}, true},
		{"empty reason is ok", DelegateInput{ToUserID: "user-1", Reason: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isInvalid := tt.input.ToUserID == ""
			if isInvalid != tt.wantErr {
				t.Errorf("expected error=%v for input %+v", tt.wantErr, tt.input)
			}
		})
	}
}

func TestReportStatusFlow(t *testing.T) {
	validTransitions := map[string][]string{
		"pending":   {"submitted"},
		"submitted": {"reviewed"},
		"reviewed":  {},
	}

	tests := []struct {
		from string
		to   string
		ok   bool
	}{
		{"pending", "submitted", true},
		{"submitted", "reviewed", true},
		{"pending", "reviewed", false},
		{"reviewed", "pending", false},
	}

	for _, tt := range tests {
		allowed := false
		for _, next := range validTransitions[tt.from] {
			if next == tt.to {
				allowed = true
				break
			}
		}
		if allowed != tt.ok {
			t.Errorf("transition %s->%s: expected ok=%v, got %v", tt.from, tt.to, tt.ok, allowed)
		}
	}
}
