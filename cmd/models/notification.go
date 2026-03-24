package models

import "time"

type Notification struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	ReferenceID *string   `json:"referenceId"`
	IsRead      bool      `json:"isRead"`
	CreatedAt   time.Time `json:"createdAt"`
}
