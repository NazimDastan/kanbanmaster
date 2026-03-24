package services

import (
	"database/sql"
	"errors"
	"time"

	"kanbanmaster/cmd/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

type CreateTaskInput struct {
	ColumnID    string  `json:"columnId"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	AssigneeID  *string `json:"assigneeId"`
	Priority    string  `json:"priority"`
	Deadline    *string `json:"deadline"`
}

func (s *TaskService) Create(creatorID string, input CreateTaskInput) (*models.Task, error) {
	// Get next position
	var maxPos sql.NullInt64
	s.db.QueryRow("SELECT MAX(position) FROM tasks WHERE column_id = $1", input.ColumnID).Scan(&maxPos)
	nextPos := 0
	if maxPos.Valid {
		nextPos = int(maxPos.Int64) + 1
	}

	if input.Priority == "" {
		input.Priority = "medium"
	}

	var deadline *time.Time
	if input.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *input.Deadline)
		if err == nil {
			deadline = &t
		}
	}

	var task models.Task
	err := s.db.QueryRow(
		`INSERT INTO tasks (column_id, title, description, creator_id, assignee_id, priority, deadline, position)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, column_id, title, description, creator_id, assignee_id, priority, deadline, position,
		           created_at, updated_at, completed_at`,
		input.ColumnID, input.Title, input.Description, creatorID,
		input.AssigneeID, input.Priority, deadline, nextPos,
	).Scan(
		&task.ID, &task.ColumnID, &task.Title, &task.Description, &task.CreatorID,
		&task.AssigneeID, &task.Priority, &task.Deadline, &task.Position,
		&task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
	)
	if err != nil {
		return nil, err
	}

	task.Subtasks = []models.Subtask{}
	task.Labels = []models.Label{}
	return &task, nil
}

func (s *TaskService) Get(taskID string) (*models.Task, error) {
	var task models.Task
	err := s.db.QueryRow(
		`SELECT id, column_id, title, description, creator_id, assignee_id, priority,
		        deadline, position, created_at, updated_at, completed_at
		 FROM tasks WHERE id = $1`,
		taskID,
	).Scan(
		&task.ID, &task.ColumnID, &task.Title, &task.Description, &task.CreatorID,
		&task.AssigneeID, &task.Priority, &task.Deadline, &task.Position,
		&task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	// Load subtasks
	subRows, err := s.db.Query(
		"SELECT id, task_id, title, is_completed, created_at FROM subtasks WHERE task_id = $1 ORDER BY created_at",
		taskID,
	)
	if err == nil {
		defer subRows.Close()
		for subRows.Next() {
			var sub models.Subtask
			subRows.Scan(&sub.ID, &sub.TaskID, &sub.Title, &sub.IsCompleted, &sub.CreatedAt)
			task.Subtasks = append(task.Subtasks, sub)
		}
	}
	if task.Subtasks == nil {
		task.Subtasks = []models.Subtask{}
	}

	task.Labels = []models.Label{}
	return &task, nil
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	AssigneeID  *string `json:"assigneeId"`
	Priority    *string `json:"priority"`
	Deadline    *string `json:"deadline"`
}

func (s *TaskService) Update(taskID string, input UpdateTaskInput) (*models.Task, error) {
	// Build dynamic update
	task, err := s.Get(taskID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = input.Description
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}

	var deadline *time.Time
	if input.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *input.Deadline)
		if err == nil {
			deadline = &t
		}
	} else {
		deadline = task.Deadline
	}

	err = s.db.QueryRow(
		`UPDATE tasks SET title=$1, description=$2, assignee_id=$3, priority=$4, deadline=$5, updated_at=NOW()
		 WHERE id=$6
		 RETURNING id, column_id, title, description, creator_id, assignee_id, priority,
		           deadline, position, created_at, updated_at, completed_at`,
		task.Title, task.Description, task.AssigneeID, task.Priority, deadline, taskID,
	).Scan(
		&task.ID, &task.ColumnID, &task.Title, &task.Description, &task.CreatorID,
		&task.AssigneeID, &task.Priority, &task.Deadline, &task.Position,
		&task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Delete(taskID string) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

type MoveTaskInput struct {
	ColumnID string `json:"columnId"`
	Position int    `json:"position"`
}

func (s *TaskService) Move(taskID string, input MoveTaskInput) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Shift tasks in target column
	_, err = tx.Exec(
		"UPDATE tasks SET position = position + 1 WHERE column_id = $1 AND position >= $2",
		input.ColumnID, input.Position,
	)
	if err != nil {
		return err
	}

	// Move task
	_, err = tx.Exec(
		"UPDATE tasks SET column_id = $1, position = $2, updated_at = NOW() WHERE id = $3",
		input.ColumnID, input.Position, taskID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TaskService) Assign(taskID string, assigneeID string) error {
	result, err := s.db.Exec(
		"UPDATE tasks SET assignee_id = $1, updated_at = NOW() WHERE id = $2",
		assigneeID, taskID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// ListByUser returns all tasks visible to a user (assigned or in their teams)
func (s *TaskService) ListByUser(userID, filter string) ([]models.Task, error) {
	baseQuery := `
		SELECT DISTINCT t.id, t.column_id, t.title, t.description, t.creator_id,
		       t.assignee_id, t.priority, t.deadline, t.position,
		       t.created_at, t.updated_at, t.completed_at,
		       c.name AS column_name
		FROM tasks t
		JOIN columns c ON c.id = t.column_id
		JOIN boards b ON b.id = c.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE tm.user_id = $1`

	switch filter {
	case "assigned":
		baseQuery += " AND t.assignee_id = $1"
	case "completed":
		baseQuery += " AND t.completed_at IS NOT NULL"
	case "overdue":
		baseQuery += " AND t.completed_at IS NULL AND t.deadline IS NOT NULL AND t.deadline < NOW()"
	case "in_progress":
		baseQuery += " AND t.completed_at IS NULL AND (t.deadline IS NULL OR t.deadline >= NOW())"
	}

	baseQuery += " ORDER BY t.updated_at DESC LIMIT 100"

	rows, err := s.db.Query(baseQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var colName string
		if err := rows.Scan(
			&t.ID, &t.ColumnID, &t.Title, &t.Description, &t.CreatorID,
			&t.AssigneeID, &t.Priority, &t.Deadline, &t.Position,
			&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &colName,
		); err != nil {
			return nil, err
		}
		t.Subtasks = []models.Subtask{}
		t.Labels = []models.Label{}
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// Subtask operations
func (s *TaskService) CreateSubtask(taskID, title string) (*models.Subtask, error) {
	var sub models.Subtask
	err := s.db.QueryRow(
		`INSERT INTO subtasks (task_id, title) VALUES ($1, $2)
		 RETURNING id, task_id, title, is_completed, created_at`,
		taskID, title,
	).Scan(&sub.ID, &sub.TaskID, &sub.Title, &sub.IsCompleted, &sub.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *TaskService) ToggleSubtask(subtaskID string) error {
	_, err := s.db.Exec(
		"UPDATE subtasks SET is_completed = NOT is_completed WHERE id = $1",
		subtaskID,
	)
	return err
}

func (s *TaskService) DeleteSubtask(subtaskID string) error {
	_, err := s.db.Exec("DELETE FROM subtasks WHERE id = $1", subtaskID)
	return err
}
