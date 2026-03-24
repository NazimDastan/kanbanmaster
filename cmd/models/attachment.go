package models

import "time"

type Attachment struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"taskId"`
	UserID      string    `json:"userId"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	Size        int       `json:"size"`
	Data        string    `json:"data"`
	CreatedAt   time.Time `json:"createdAt"`
}
