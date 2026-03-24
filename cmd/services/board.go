package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var (
	ErrBoardNotFound  = errors.New("board not found")
	ErrColumnNotFound = errors.New("column not found")
)

type BoardService struct {
	db *sql.DB
}

func NewBoardService(db *sql.DB) *BoardService {
	return &BoardService{db: db}
}

func (s *BoardService) Create(name, teamID string) (*models.Board, error) {
	var board models.Board
	err := s.db.QueryRow(
		`INSERT INTO boards (name, team_id) VALUES ($1, $2)
		 RETURNING id, name, team_id, created_at`,
		name, teamID,
	).Scan(&board.ID, &board.Name, &board.TeamID, &board.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Create default columns
	defaults := []string{"Todo", "In Progress", "Done"}
	for i, colName := range defaults {
		_, err := s.db.Exec(
			"INSERT INTO columns (board_id, name, position) VALUES ($1, $2, $3)",
			board.ID, colName, i,
		)
		if err != nil {
			return nil, err
		}
	}

	return &board, nil
}

func (s *BoardService) List(userID string) ([]models.Board, error) {
	rows, err := s.db.Query(
		`SELECT b.id, b.name, b.team_id, b.created_at
		 FROM boards b
		 JOIN teams t ON t.id = b.team_id
		 JOIN team_members tm ON tm.team_id = t.id
		 WHERE tm.user_id = $1
		 ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var b models.Board
		if err := rows.Scan(&b.ID, &b.Name, &b.TeamID, &b.CreatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (s *BoardService) GetWithColumns(boardID string) (*models.BoardWithColumns, error) {
	var board models.BoardWithColumns
	err := s.db.QueryRow(
		"SELECT id, name, team_id, created_at FROM boards WHERE id = $1",
		boardID,
	).Scan(&board.ID, &board.Name, &board.TeamID, &board.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}

	// Fetch columns
	colRows, err := s.db.Query(
		"SELECT id, board_id, name, position, created_at FROM columns WHERE board_id = $1 ORDER BY position",
		boardID,
	)
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	for colRows.Next() {
		var col models.Column
		if err := colRows.Scan(&col.ID, &col.BoardID, &col.Name, &col.Position, &col.CreatedAt); err != nil {
			return nil, err
		}
		col.Tasks = []models.Task{}
		board.Columns = append(board.Columns, col)
	}

	// Fetch tasks for each column
	for i := range board.Columns {
		tasks, err := s.getColumnTasks(board.Columns[i].ID)
		if err != nil {
			return nil, err
		}
		board.Columns[i].Tasks = tasks
	}

	if board.Columns == nil {
		board.Columns = []models.Column{}
	}

	return &board, nil
}

func (s *BoardService) Update(boardID, name string) (*models.Board, error) {
	var board models.Board
	err := s.db.QueryRow(
		`UPDATE boards SET name = $1 WHERE id = $2
		 RETURNING id, name, team_id, created_at`,
		name, boardID,
	).Scan(&board.ID, &board.Name, &board.TeamID, &board.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}
	return &board, nil
}

func (s *BoardService) Delete(boardID string) error {
	result, err := s.db.Exec("DELETE FROM boards WHERE id = $1", boardID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrBoardNotFound
	}
	return nil
}

func (s *BoardService) getColumnTasks(columnID string) ([]models.Task, error) {
	rows, err := s.db.Query(
		`SELECT t.id, t.column_id, t.title, t.description, t.creator_id,
		        t.assignee_id, t.priority, t.deadline, t.position,
		        t.created_at, t.updated_at, t.completed_at
		 FROM tasks t
		 WHERE t.column_id = $1
		 ORDER BY t.position`,
		columnID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID, &t.ColumnID, &t.Title, &t.Description, &t.CreatorID,
			&t.AssigneeID, &t.Priority, &t.Deadline, &t.Position,
			&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
		); err != nil {
			return nil, err
		}

		// Load assignee if exists
		if t.AssigneeID != nil {
			var user models.User
			err := s.db.QueryRow(
				"SELECT id, email, name, avatar_url, created_at, updated_at FROM users WHERE id = $1",
				*t.AssigneeID,
			).Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
			if err == nil {
				t.Assignee = &user
			}
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
