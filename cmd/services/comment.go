package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentService struct {
	db *sql.DB
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{db: db}
}

type Comment struct {
	ID        string      `json:"id"`
	TaskID    string      `json:"taskId"`
	UserID    string      `json:"userId"`
	User      *models.User `json:"user,omitempty"`
	Content   string      `json:"content"`
	CreatedAt string      `json:"createdAt"`
}

func (s *CommentService) Create(taskID, userID, content string) (*Comment, error) {
	var c Comment
	err := s.db.QueryRow(
		`INSERT INTO comments (task_id, user_id, content) VALUES ($1, $2, $3)
		 RETURNING id, task_id, user_id, content, created_at`,
		taskID, userID, content,
	).Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Load user
	var user models.User
	err = s.db.QueryRow(
		"SELECT id, email, name, avatar_url, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
	if err == nil {
		c.User = &user
	}

	return &c, nil
}

func (s *CommentService) ListByTask(taskID string) ([]Comment, error) {
	rows, err := s.db.Query(
		`SELECT c.id, c.task_id, c.user_id, c.content, c.created_at,
		        u.id, u.email, u.name, u.avatar_url, u.created_at, u.updated_at
		 FROM comments c
		 JOIN users u ON u.id = c.user_id
		 WHERE c.task_id = $1
		 ORDER BY c.created_at ASC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var u models.User
		if err := rows.Scan(
			&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt,
			&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.User = &u
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []Comment{}
	}
	return comments, nil
}

func (s *CommentService) Delete(commentID, userID string) error {
	result, err := s.db.Exec(
		"DELETE FROM comments WHERE id = $1 AND user_id = $2",
		commentID, userID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("comment not found or not authorized")
	}
	return nil
}
