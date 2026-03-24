package models

import "time"

type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TeamID    string    `json:"teamId"`
	CreatedAt time.Time `json:"createdAt"`
}

type BoardWithColumns struct {
	Board
	Columns []Column `json:"columns"`
}
