package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var ErrLabelNotFound = errors.New("label not found")

type LabelService struct {
	db *sql.DB
}

func NewLabelService(db *sql.DB) *LabelService {
	return &LabelService{db: db}
}

func (s *LabelService) Create(boardID, name, color string) (*models.Label, error) {
	if color == "" {
		color = "#1E88E5"
	}
	var label models.Label
	err := s.db.QueryRow(
		`INSERT INTO labels (board_id, name, color) VALUES ($1, $2, $3)
		 RETURNING id, board_id, name, color`,
		boardID, name, color,
	).Scan(&label.ID, &label.BoardID, &label.Name, &label.Color)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (s *LabelService) ListByBoard(boardID string) ([]models.Label, error) {
	rows, err := s.db.Query(
		"SELECT id, board_id, name, color FROM labels WHERE board_id = $1 ORDER BY name",
		boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var l models.Label
		if err := rows.Scan(&l.ID, &l.BoardID, &l.Name, &l.Color); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	if labels == nil {
		labels = []models.Label{}
	}
	return labels, nil
}

func (s *LabelService) Update(labelID, name, color string) (*models.Label, error) {
	var label models.Label
	err := s.db.QueryRow(
		`UPDATE labels SET name = $1, color = $2 WHERE id = $3
		 RETURNING id, board_id, name, color`,
		name, color, labelID,
	).Scan(&label.ID, &label.BoardID, &label.Name, &label.Color)
	if err != nil {
		return nil, ErrLabelNotFound
	}
	return &label, nil
}

func (s *LabelService) Delete(labelID string) error {
	_, err := s.db.Exec("DELETE FROM labels WHERE id = $1", labelID)
	return err
}

func (s *LabelService) AddToTask(taskID, labelID string) error {
	_, err := s.db.Exec(
		"INSERT INTO task_labels (task_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		taskID, labelID,
	)
	return err
}

func (s *LabelService) RemoveFromTask(taskID, labelID string) error {
	_, err := s.db.Exec(
		"DELETE FROM task_labels WHERE task_id = $1 AND label_id = $2",
		taskID, labelID,
	)
	return err
}

func (s *LabelService) GetTaskLabels(taskID string) ([]models.Label, error) {
	rows, err := s.db.Query(
		`SELECT l.id, l.board_id, l.name, l.color
		 FROM labels l
		 JOIN task_labels tl ON tl.label_id = l.id
		 WHERE tl.task_id = $1`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var l models.Label
		if err := rows.Scan(&l.ID, &l.BoardID, &l.Name, &l.Color); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	if labels == nil {
		labels = []models.Label{}
	}
	return labels, nil
}
