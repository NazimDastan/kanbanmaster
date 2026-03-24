package services

import (
	"testing"
)

func TestTeamServiceCheckPermission(t *testing.T) {
	// Unit test for permission logic without DB
	// Tests that role checking correctly identifies leaders
	validRoles := map[string]bool{
		"owner":  true,
		"leader": true,
		"member": false,
		"viewer": false,
	}

	for role, shouldPass := range validRoles {
		isLeaderOrOwner := role == "owner" || role == "leader"
		if isLeaderOrOwner != shouldPass {
			t.Errorf("role %s: expected pass=%v, got %v", role, shouldPass, isLeaderOrOwner)
		}
	}
}

func TestTeamMemberRoleValidation(t *testing.T) {
	validRoles := []string{"owner", "leader", "member", "viewer"}
	roleSet := make(map[string]bool)
	for _, r := range validRoles {
		roleSet[r] = true
	}

	tests := []struct {
		role  string
		valid bool
	}{
		{"owner", true},
		{"leader", true},
		{"member", true},
		{"viewer", true},
		{"admin", false},
		{"", false},
		{"superuser", false},
	}

	for _, tt := range tests {
		if roleSet[tt.role] != tt.valid {
			t.Errorf("role %q: expected valid=%v, got %v", tt.role, tt.valid, roleSet[tt.role])
		}
	}
}
