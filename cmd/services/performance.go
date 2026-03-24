package services

import (
	"database/sql"
	"time"
)

type PerformanceService struct {
	db *sql.DB
}

func NewPerformanceService(db *sql.DB) *PerformanceService {
	return &PerformanceService{db: db}
}

type DashboardSummary struct {
	TotalTasks   int `json:"totalTasks"`
	Completed    int `json:"completed"`
	InProgress   int `json:"inProgress"`
	Overdue      int `json:"overdue"`
	CompletedPct int `json:"completedPct"`
}

func (s *PerformanceService) GetDashboardSummary(userID string) (*DashboardSummary, error) {
	var summary DashboardSummary

	s.db.QueryRow(
		`SELECT COUNT(*) FROM tasks t
		 JOIN columns c ON c.id = t.column_id
		 JOIN boards b ON b.id = c.board_id
		 JOIN team_members tm ON tm.team_id = b.team_id
		 WHERE tm.user_id = $1`, userID,
	).Scan(&summary.TotalTasks)

	s.db.QueryRow(
		`SELECT COUNT(*) FROM tasks t
		 JOIN columns c ON c.id = t.column_id
		 JOIN boards b ON b.id = c.board_id
		 JOIN team_members tm ON tm.team_id = b.team_id
		 WHERE tm.user_id = $1 AND t.completed_at IS NOT NULL`, userID,
	).Scan(&summary.Completed)

	s.db.QueryRow(
		`SELECT COUNT(*) FROM tasks t
		 JOIN columns c ON c.id = t.column_id
		 JOIN boards b ON b.id = c.board_id
		 JOIN team_members tm ON tm.team_id = b.team_id
		 WHERE tm.user_id = $1 AND t.completed_at IS NULL AND t.deadline < NOW()`, userID,
	).Scan(&summary.Overdue)

	summary.InProgress = summary.TotalTasks - summary.Completed - summary.Overdue
	if summary.InProgress < 0 {
		summary.InProgress = 0
	}

	if summary.TotalTasks > 0 {
		summary.CompletedPct = (summary.Completed * 100) / summary.TotalTasks
	}

	return &summary, nil
}

type UserPerformance struct {
	UserID       string `json:"userId"`
	UserName     string `json:"userName"`
	TotalTasks   int    `json:"totalTasks"`
	Completed    int    `json:"completed"`
	OnTime       int    `json:"onTime"`
	Overdue      int    `json:"overdue"`
	Score        int    `json:"score"`
}

func (s *PerformanceService) GetTeamPerformance(teamID string) ([]UserPerformance, error) {
	rows, err := s.db.Query(
		`SELECT u.id, u.name,
		        COUNT(t.id) AS total,
		        COUNT(t.completed_at) AS completed,
		        COUNT(CASE WHEN t.completed_at IS NOT NULL AND (t.deadline IS NULL OR t.completed_at <= t.deadline) THEN 1 END) AS on_time,
		        COUNT(CASE WHEN t.completed_at IS NULL AND t.deadline IS NOT NULL AND t.deadline < NOW() THEN 1 END) AS overdue
		 FROM team_members tm
		 JOIN users u ON u.id = tm.user_id
		 LEFT JOIN tasks t ON t.assignee_id = u.id
		 WHERE tm.team_id = $1
		 GROUP BY u.id, u.name
		 ORDER BY COUNT(t.completed_at) DESC`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []UserPerformance
	for rows.Next() {
		var p UserPerformance
		if err := rows.Scan(&p.UserID, &p.UserName, &p.TotalTasks, &p.Completed, &p.OnTime, &p.Overdue); err != nil {
			return nil, err
		}
		if p.TotalTasks > 0 {
			p.Score = (p.OnTime * 100) / p.TotalTasks
		}
		results = append(results, p)
	}
	if results == nil {
		results = []UserPerformance{}
	}
	return results, nil
}

type OverdueTask struct {
	TaskID    string    `json:"taskId"`
	Title     string    `json:"title"`
	Assignee  string    `json:"assignee"`
	Deadline  time.Time `json:"deadline"`
	DaysLate  int       `json:"daysLate"`
}

func (s *PerformanceService) GetOverdueTasks(userID string) ([]OverdueTask, error) {
	rows, err := s.db.Query(
		`SELECT t.id, t.title, COALESCE(u.name, 'Unassigned'), t.deadline
		 FROM tasks t
		 JOIN columns c ON c.id = t.column_id
		 JOIN boards b ON b.id = c.board_id
		 JOIN team_members tm ON tm.team_id = b.team_id
		 LEFT JOIN users u ON u.id = t.assignee_id
		 WHERE tm.user_id = $1 AND t.completed_at IS NULL AND t.deadline < NOW()
		 ORDER BY t.deadline ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []OverdueTask
	now := time.Now()
	for rows.Next() {
		var t OverdueTask
		if err := rows.Scan(&t.TaskID, &t.Title, &t.Assignee, &t.Deadline); err != nil {
			return nil, err
		}
		t.DaysLate = int(now.Sub(t.Deadline).Hours() / 24)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []OverdueTask{}
	}
	return tasks, nil
}
