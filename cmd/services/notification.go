package services

import (
	"database/sql"

	"kanbanmaster/cmd/models"
)

// OnNotifyFunc is a callback invoked after a notification is created
type OnNotifyFunc func(userID string, n models.Notification)

type NotificationService struct {
	db       *sql.DB
	onNotify OnNotifyFunc
}

// SetOnNotify sets the real-time broadcast callback
func (s *NotificationService) SetOnNotify(fn OnNotifyFunc) {
	s.onNotify = fn
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) Create(userID, nType, title, message string, referenceID *string) (*models.Notification, error) {
	var n models.Notification
	err := s.db.QueryRow(
		`INSERT INTO notifications (user_id, type, title, message, reference_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, type, title, message, reference_id, is_read, created_at`,
		userID, nType, title, message, referenceID,
	).Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.ReferenceID, &n.IsRead, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	// Broadcast via WebSocket if callback set
	if s.onNotify != nil {
		s.onNotify(userID, n)
	}
	return &n, nil
}

func (s *NotificationService) ListByUser(userID string) ([]models.Notification, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, type, title, message, reference_id, is_read, created_at
		 FROM notifications WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.ReferenceID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	if notifications == nil {
		notifications = []models.Notification{}
	}
	return notifications, nil
}

func (s *NotificationService) MarkRead(notifID, userID string) error {
	_, err := s.db.Exec(
		"UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2",
		notifID, userID,
	)
	return err
}

func (s *NotificationService) MarkAllRead(userID string) error {
	_, err := s.db.Exec(
		"UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false",
		userID,
	)
	return err
}

func (s *NotificationService) UnreadCount(userID string) (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false",
		userID,
	).Scan(&count)
	return count, err
}

// Notify helpers — called by other services after events
func (s *NotificationService) NotifyTaskAssigned(assigneeID, taskTitle, taskID string) {
	s.Create(assigneeID, "assigned", "Task Assigned", "You were assigned to: "+taskTitle, &taskID)
}

func (s *NotificationService) NotifyTaskDelegated(toUserID, fromUserName, taskTitle, taskID string) {
	s.Create(toUserID, "delegated", "Task Delegated", fromUserName+" delegated: "+taskTitle, &taskID)
}

func (s *NotificationService) NotifyDeadlineApproaching(userID, taskTitle, taskID string) {
	s.Create(userID, "deadline", "Deadline Approaching", taskTitle+" is due soon", &taskID)
}

func (s *NotificationService) NotifyTaskCompleted(leaderID, userName, taskTitle, taskID string) {
	s.Create(leaderID, "completed", "Task Completed", userName+" completed: "+taskTitle, &taskID)
}

func (s *NotificationService) NotifyComment(userID, commenterName, taskTitle, taskID string) {
	s.Create(userID, "comment", "New Comment", commenterName+" commented on: "+taskTitle, &taskID)
}

func (s *NotificationService) NotifyReportRequested(targetUserID, requesterName, reportID string) {
	s.Create(targetUserID, "report_request", "Report Requested", requesterName+" requested a report from you", &reportID)
}

func (s *NotificationService) NotifyOverdue(userID, taskTitle, taskID string) {
	s.Create(userID, "overdue", "Task Overdue", taskTitle+" is past its deadline", &taskID)
}
