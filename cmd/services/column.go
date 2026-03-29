package services

import (
	"database/sql"

	"kanbanmaster/cmd/models"
)

type ColumnService struct {
	db *sql.DB
}

func NewColumnService(db *sql.DB) *ColumnService {
	return &ColumnService{db: db}
}

func (s *ColumnService) Create(boardID, name, color string) (*models.Column, error) {
	// Get next position
	var maxPos sql.NullInt64
	s.db.QueryRow("SELECT MAX(position) FROM columns WHERE board_id = $1", boardID).Scan(&maxPos)
	nextPos := 0
	if maxPos.Valid {
		nextPos = int(maxPos.Int64) + 1
	}

	var col models.Column
	err := s.db.QueryRow(
		`INSERT INTO columns (board_id, name, position, color) VALUES ($1, $2, $3, $4)
		 RETURNING id, board_id, name, position, color, created_at`,
		boardID, name, nextPos, color,
	).Scan(&col.ID, &col.BoardID, &col.Name, &col.Position, &col.Color, &col.CreatedAt)
	if err != nil {
		return nil, err
	}
	col.Tasks = []models.Task{}
	return &col, nil
}

func (s *ColumnService) Update(columnID, name string) (*models.Column, error) {
	var col models.Column
	err := s.db.QueryRow(
		`UPDATE columns SET name = $1 WHERE id = $2
		 RETURNING id, board_id, name, position, color, created_at`,
		name, columnID,
	).Scan(&col.ID, &col.BoardID, &col.Name, &col.Position, &col.Color, &col.CreatedAt)
	if err != nil {
		return nil, ErrColumnNotFound
	}
	return &col, nil
}

func (s *ColumnService) Delete(columnID string) error {
	_, err := s.db.Exec("DELETE FROM columns WHERE id = $1", columnID)
	return err
}

type ReorderInput struct {
	ColumnID string `json:"columnId"`
	Position int    `json:"position"`
}

func (s *ColumnService) Reorder(boardID string, items []ReorderInput) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.Exec(
			"UPDATE columns SET position = $1 WHERE id = $2 AND board_id = $3",
			item.Position, item.ColumnID, boardID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
