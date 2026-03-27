package models

import "time"

type Invitation struct {
	ID          string     `json:"id"`
	TeamID      string     `json:"teamId"`
	Team        *Team      `json:"team,omitempty"`
	InviterID   string     `json:"inviterId"`
	Inviter     *User      `json:"inviter,omitempty"`
	InviteeID   string     `json:"inviteeId"`
	Invitee     *User      `json:"invitee,omitempty"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	RespondedAt *time.Time `json:"respondedAt"`
}
