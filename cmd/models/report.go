package models

import "time"

type TaskDelegation struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"taskId"`
	FromUserID  string    `json:"fromUserId"`
	FromUser    *User     `json:"fromUser,omitempty"`
	ToUserID    string    `json:"toUserId"`
	ToUser      *User     `json:"toUser,omitempty"`
	Reason      string    `json:"reason"`
	DelegatedAt time.Time `json:"delegatedAt"`
}

type ReportRequest struct {
	ID           string     `json:"id"`
	RequesterID  string     `json:"requesterId"`
	Requester    *User      `json:"requester,omitempty"`
	TargetUserID string     `json:"targetUserId"`
	TargetUser   *User      `json:"targetUser,omitempty"`
	TeamID       string     `json:"teamId"`
	Message      string     `json:"message"`
	Response     *string    `json:"response"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	RespondedAt  *time.Time `json:"respondedAt"`
}

type ActivityLog struct {
	ID        string                 `json:"id"`
	TaskID    string                 `json:"taskId"`
	UserID    string                 `json:"userId"`
	User      *User                  `json:"user,omitempty"`
	Action    string                 `json:"action"`
	Details   map[string]interface{} `json:"details"`
	CreatedAt time.Time              `json:"createdAt"`
}
