package services

import (
	"database/sql"
	"log"
	"time"
)

type Scheduler struct {
	db              *sql.DB
	notifService    *NotificationService
}

func NewScheduler(db *sql.DB, notifService *NotificationService) *Scheduler {
	return &Scheduler{db: db, notifService: notifService}
}

// Start runs the scheduler in a goroutine, checking every hour
func (s *Scheduler) Start() {
	go func() {
		// Run immediately on start
		s.checkDeadlines()

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			s.checkDeadlines()
		}
	}()
	log.Println("Deadline scheduler started (checking every hour)")
}

func (s *Scheduler) checkDeadlines() {
	// 1. Tasks due within 24 hours (not yet notified recently)
	s.notifyApproachingDeadlines()

	// 2. Tasks that are overdue (not yet completed)
	s.notifyOverdueTasks()
}

func (s *Scheduler) notifyApproachingDeadlines() {
	rows, err := s.db.Query(`
		SELECT t.id, t.title, t.assignee_id, t.deadline
		FROM tasks t
		WHERE t.assignee_id IS NOT NULL
		  AND t.completed_at IS NULL
		  AND t.deadline IS NOT NULL
		  AND t.deadline BETWEEN NOW() AND NOW() + INTERVAL '24 hours'
		  AND NOT EXISTS (
		    SELECT 1 FROM notifications n
		    WHERE n.reference_id::text = t.id::text
		      AND n.type = 'deadline'
		      AND n.created_at > NOW() - INTERVAL '24 hours'
		  )
	`)
	if err != nil {
		log.Printf("Scheduler: deadline query error: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var taskID, title, assigneeID string
		var deadline time.Time
		if err := rows.Scan(&taskID, &title, &assigneeID, &deadline); err != nil {
			continue
		}
		s.notifService.NotifyDeadlineApproaching(assigneeID, title, taskID)
		count++
	}
	if count > 0 {
		log.Printf("Scheduler: sent %d deadline approaching notifications", count)
	}
}

func (s *Scheduler) notifyOverdueTasks() {
	rows, err := s.db.Query(`
		SELECT t.id, t.title, t.assignee_id
		FROM tasks t
		WHERE t.assignee_id IS NOT NULL
		  AND t.completed_at IS NULL
		  AND t.deadline IS NOT NULL
		  AND t.deadline < NOW()
		  AND NOT EXISTS (
		    SELECT 1 FROM notifications n
		    WHERE n.reference_id::text = t.id::text
		      AND n.type = 'overdue'
		      AND n.created_at > NOW() - INTERVAL '24 hours'
		  )
	`)
	if err != nil {
		log.Printf("Scheduler: overdue query error: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var taskID, title, assigneeID string
		if err := rows.Scan(&taskID, &title, &assigneeID); err != nil {
			continue
		}
		s.notifService.NotifyOverdue(assigneeID, title, taskID)
		count++
	}
	if count > 0 {
		log.Printf("Scheduler: sent %d overdue notifications", count)
	}
}
