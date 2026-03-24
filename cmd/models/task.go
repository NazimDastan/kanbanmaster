package models

import "time"

type Task struct {
	ID          string     `json:"id"`
	ColumnID    string     `json:"columnId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	CreatorID   string     `json:"creatorId"`
	AssigneeID  *string    `json:"assigneeId"`
	Assignee    *User      `json:"assignee,omitempty"`
	Priority    string     `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	Position    int        `json:"position"`
	Subtasks    []Subtask  `json:"subtasks"`
	Labels      []Label    `json:"labels"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type Subtask struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"taskId"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Label struct {
	ID      string `json:"id"`
	BoardID string `json:"boardId"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}
