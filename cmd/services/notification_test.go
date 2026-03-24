package services

import (
	"testing"
)

func TestNotificationTypes(t *testing.T) {
	validTypes := map[string]bool{
		"assigned":       true,
		"delegated":      true,
		"deadline":       true,
		"comment":        true,
		"report_request": true,
		"completed":      true,
		"overdue":        true,
	}

	tests := []string{"assigned", "delegated", "deadline", "comment", "report_request", "completed", "overdue", "invalid", ""}
	for _, nType := range tests {
		if validTypes[nType] && nType == "" {
			t.Error("empty type should not be valid")
		}
		if !validTypes[nType] && nType != "invalid" && nType != "" {
			t.Errorf("type %q should be valid", nType)
		}
	}
}
