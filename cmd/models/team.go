package models

import "time"

type Team struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	OrganizationID string       `json:"organizationId"`
	CreatedAt      time.Time    `json:"createdAt"`
	Members        []TeamMember `json:"members,omitempty"`
}

type TeamMember struct {
	ID       string    `json:"id"`
	TeamID   string    `json:"teamId"`
	UserID   string    `json:"userId"`
	User     *User     `json:"user,omitempty"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
}
