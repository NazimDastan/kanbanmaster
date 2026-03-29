package models

import "time"

type Column struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"boardId"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	Color     string    `json:"color"`
	Tasks     []Task    `json:"tasks"`
	CreatedAt time.Time `json:"createdAt"`
}
