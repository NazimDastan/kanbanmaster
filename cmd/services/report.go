package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var (
	ErrReportNotFound = errors.New("report request not found")
	ErrNotAuthorized  = errors.New("not authorized for this action")
)

type ReportService struct {
	db *sql.DB
}

func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{db: db}
}

type CreateReportRequest struct {
	TargetUserID string `json:"targetUserId"`
	TeamID       string `json:"teamId"`
	Message      string `json:"message"`
}

func (s *ReportService) RequestReport(requesterID string, input CreateReportRequest) (*models.ReportRequest, error) {
	var rr models.ReportRequest
	err := s.db.QueryRow(
		`INSERT INTO report_requests (requester_id, target_user_id, team_id, message)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, requester_id, target_user_id, team_id, message, response, status, created_at, responded_at`,
		requesterID, input.TargetUserID, input.TeamID, input.Message,
	).Scan(&rr.ID, &rr.RequesterID, &rr.TargetUserID, &rr.TeamID,
		&rr.Message, &rr.Response, &rr.Status, &rr.CreatedAt, &rr.RespondedAt)
	if err != nil {
		return nil, err
	}
	return &rr, nil
}

func (s *ReportService) GetIncoming(userID string) ([]models.ReportRequest, error) {
	return s.queryReports(
		`SELECT r.id, r.requester_id, r.target_user_id, r.team_id, r.message,
		        r.response, r.status, r.created_at, r.responded_at,
		        u.id, u.name, u.email
		 FROM report_requests r
		 JOIN users u ON u.id = r.requester_id
		 WHERE r.target_user_id = $1
		 ORDER BY r.created_at DESC`, userID, true,
	)
}

func (s *ReportService) GetSent(userID string) ([]models.ReportRequest, error) {
	return s.queryReports(
		`SELECT r.id, r.requester_id, r.target_user_id, r.team_id, r.message,
		        r.response, r.status, r.created_at, r.responded_at,
		        u.id, u.name, u.email
		 FROM report_requests r
		 JOIN users u ON u.id = r.target_user_id
		 WHERE r.requester_id = $1
		 ORDER BY r.created_at DESC`, userID, false,
	)
}

func (s *ReportService) Respond(reportID, userID, response string) (*models.ReportRequest, error) {
	var rr models.ReportRequest
	err := s.db.QueryRow(
		`UPDATE report_requests
		 SET response = $1, status = 'submitted', responded_at = NOW()
		 WHERE id = $2 AND target_user_id = $3
		 RETURNING id, requester_id, target_user_id, team_id, message, response, status, created_at, responded_at`,
		response, reportID, userID,
	).Scan(&rr.ID, &rr.RequesterID, &rr.TargetUserID, &rr.TeamID,
		&rr.Message, &rr.Response, &rr.Status, &rr.CreatedAt, &rr.RespondedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotAuthorized
		}
		return nil, err
	}
	return &rr, nil
}

func (s *ReportService) Review(reportID, userID string) (*models.ReportRequest, error) {
	var rr models.ReportRequest
	err := s.db.QueryRow(
		`UPDATE report_requests SET status = 'reviewed'
		 WHERE id = $1 AND requester_id = $2
		 RETURNING id, requester_id, target_user_id, team_id, message, response, status, created_at, responded_at`,
		reportID, userID,
	).Scan(&rr.ID, &rr.RequesterID, &rr.TargetUserID, &rr.TeamID,
		&rr.Message, &rr.Response, &rr.Status, &rr.CreatedAt, &rr.RespondedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotAuthorized
		}
		return nil, err
	}
	return &rr, nil
}

func (s *ReportService) queryReports(query, userID string, isIncoming bool) ([]models.ReportRequest, error) {
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.ReportRequest
	for rows.Next() {
		var r models.ReportRequest
		var otherUser models.User
		err := rows.Scan(
			&r.ID, &r.RequesterID, &r.TargetUserID, &r.TeamID,
			&r.Message, &r.Response, &r.Status, &r.CreatedAt, &r.RespondedAt,
			&otherUser.ID, &otherUser.Name, &otherUser.Email,
		)
		if err != nil {
			return nil, err
		}
		if isIncoming {
			r.Requester = &otherUser
		} else {
			r.TargetUser = &otherUser
		}
		reports = append(reports, r)
	}
	if reports == nil {
		reports = []models.ReportRequest{}
	}
	return reports, nil
}
