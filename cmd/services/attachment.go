package services

import (
	"database/sql"

	"kanbanmaster/cmd/models"
)

type AttachmentService struct {
	db *sql.DB
}

func NewAttachmentService(db *sql.DB) *AttachmentService {
	return &AttachmentService{db: db}
}

func (s *AttachmentService) Create(taskID, userID, filename, contentType string, size int, data string) (*models.Attachment, error) {
	var a models.Attachment
	err := s.db.QueryRow(
		`INSERT INTO task_attachments (task_id, user_id, filename, content_type, size, data)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, task_id, user_id, filename, content_type, size, data, created_at`,
		taskID, userID, filename, contentType, size, data,
	).Scan(&a.ID, &a.TaskID, &a.UserID, &a.Filename, &a.ContentType, &a.Size, &a.Data, &a.CreatedAt)
	return &a, err
}

func (s *AttachmentService) ListByTask(taskID string) ([]models.Attachment, error) {
	rows, err := s.db.Query(
		`SELECT id, task_id, user_id, filename, content_type, size, '' as data, created_at
		 FROM task_attachments WHERE task_id = $1 ORDER BY created_at DESC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.TaskID, &a.UserID, &a.Filename, &a.ContentType, &a.Size, &a.Data, &a.CreatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	if attachments == nil {
		attachments = []models.Attachment{}
	}
	return attachments, nil
}

func (s *AttachmentService) Get(id string) (*models.Attachment, error) {
	var a models.Attachment
	err := s.db.QueryRow(
		"SELECT id, task_id, user_id, filename, content_type, size, data, created_at FROM task_attachments WHERE id = $1",
		id,
	).Scan(&a.ID, &a.TaskID, &a.UserID, &a.Filename, &a.ContentType, &a.Size, &a.Data, &a.CreatedAt)
	return &a, err
}

func (s *AttachmentService) Delete(id, userID string) error {
	_, err := s.db.Exec("DELETE FROM task_attachments WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
