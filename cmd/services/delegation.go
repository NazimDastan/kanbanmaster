package services

import (
	"database/sql"
	"encoding/json"
	"errors"

	"kanbanmaster/cmd/models"
)

var ErrDelegationNotFound = errors.New("delegation not found")

type DelegationService struct {
	db *sql.DB
}

func NewDelegationService(db *sql.DB) *DelegationService {
	return &DelegationService{db: db}
}

type DelegateInput struct {
	ToUserID string `json:"toUserId"`
	Reason   string `json:"reason"`
}

func (s *DelegationService) Delegate(taskID, fromUserID string, input DelegateInput) (*models.TaskDelegation, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create delegation record
	var d models.TaskDelegation
	err = tx.QueryRow(
		`INSERT INTO task_delegations (task_id, from_user_id, to_user_id, reason)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, task_id, from_user_id, to_user_id, reason, delegated_at`,
		taskID, fromUserID, input.ToUserID, input.Reason,
	).Scan(&d.ID, &d.TaskID, &d.FromUserID, &d.ToUserID, &d.Reason, &d.DelegatedAt)
	if err != nil {
		return nil, err
	}

	// Update task assignee
	_, err = tx.Exec(
		"UPDATE tasks SET assignee_id = $1, updated_at = NOW() WHERE id = $2",
		input.ToUserID, taskID,
	)
	if err != nil {
		return nil, err
	}

	// Log activity (safe JSON encoding)
	detailsJSON, _ := json.Marshal(map[string]string{
		"fromUserId": fromUserID,
		"toUserId":   input.ToUserID,
		"reason":     input.Reason,
	})
	_, err = tx.Exec(
		`INSERT INTO task_activity_log (task_id, user_id, action, details)
		 VALUES ($1, $2, 'delegated', $3)`,
		taskID, fromUserID, string(detailsJSON),
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *DelegationService) GetTaskDelegations(taskID string) ([]models.TaskDelegation, error) {
	rows, err := s.db.Query(
		`SELECT d.id, d.task_id, d.from_user_id, d.to_user_id, d.reason, d.delegated_at,
		        fu.id, fu.name, fu.email, tu.id, tu.name, tu.email
		 FROM task_delegations d
		 JOIN users fu ON fu.id = d.from_user_id
		 JOIN users tu ON tu.id = d.to_user_id
		 WHERE d.task_id = $1
		 ORDER BY d.delegated_at DESC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var delegations []models.TaskDelegation
	for rows.Next() {
		var d models.TaskDelegation
		var fromUser, toUser models.User
		err := rows.Scan(
			&d.ID, &d.TaskID, &d.FromUserID, &d.ToUserID, &d.Reason, &d.DelegatedAt,
			&fromUser.ID, &fromUser.Name, &fromUser.Email,
			&toUser.ID, &toUser.Name, &toUser.Email,
		)
		if err != nil {
			return nil, err
		}
		d.FromUser = &fromUser
		d.ToUser = &toUser
		delegations = append(delegations, d)
	}
	if delegations == nil {
		delegations = []models.TaskDelegation{}
	}
	return delegations, nil
}

func (s *DelegationService) GetActivityLog(taskID string) ([]models.ActivityLog, error) {
	rows, err := s.db.Query(
		`SELECT a.id, a.task_id, a.user_id, a.action, a.details, a.created_at,
		        u.id, u.name, u.email
		 FROM task_activity_log a
		 JOIN users u ON u.id = a.user_id
		 WHERE a.task_id = $1
		 ORDER BY a.created_at DESC`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		var user models.User
		var detailsJSON []byte
		err := rows.Scan(
			&log.ID, &log.TaskID, &log.UserID, &log.Action, &detailsJSON, &log.CreatedAt,
			&user.ID, &user.Name, &user.Email,
		)
		if err != nil {
			return nil, err
		}
		log.User = &user
		log.Details = make(map[string]interface{})
		logs = append(logs, log)
	}
	if logs == nil {
		logs = []models.ActivityLog{}
	}
	return logs, nil
}

// LogActivity records an action on a task
func (s *DelegationService) LogActivity(taskID, userID, action, details string) error {
	_, err := s.db.Exec(
		`INSERT INTO task_activity_log (task_id, user_id, action, details) VALUES ($1, $2, $3, $4)`,
		taskID, userID, action, details,
	)
	return err
}
