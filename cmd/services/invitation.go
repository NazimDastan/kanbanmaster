package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var (
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrAlreadyInvited     = errors.New("user already invited")
	ErrNotInvitee         = errors.New("not the invitee")
)

type InvitationService struct {
	db           *sql.DB
	notifService *NotificationService
}

func NewInvitationService(db *sql.DB, notifService *NotificationService) *InvitationService {
	return &InvitationService{db: db, notifService: notifService}
}

// Send creates a pending invitation
func (s *InvitationService) Send(teamID, inviterID, inviteeEmail, role string) (*models.Invitation, error) {
	// Find invitee by email
	var inviteeID string
	err := s.db.QueryRow("SELECT id FROM users WHERE email = $1", inviteeEmail).Scan(&inviteeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check not already member
	var isMember bool
	s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)", teamID, inviteeID).Scan(&isMember)
	if isMember {
		return nil, ErrAlreadyMember
	}

	// Check no pending invitation
	var hasPending bool
	s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM team_invitations WHERE team_id = $1 AND invitee_id = $2 AND status = 'pending')", teamID, inviteeID).Scan(&hasPending)
	if hasPending {
		return nil, ErrAlreadyInvited
	}

	if role == "" {
		role = "member"
	}

	var inv models.Invitation
	err = s.db.QueryRow(
		`INSERT INTO team_invitations (team_id, inviter_id, invitee_id, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, team_id, inviter_id, invitee_id, role, status, created_at, responded_at`,
		teamID, inviterID, inviteeID, role,
	).Scan(&inv.ID, &inv.TeamID, &inv.InviterID, &inv.InviteeID, &inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt)
	if err != nil {
		return nil, err
	}

	// Send notification to invitee
	var inviterName string
	s.db.QueryRow("SELECT name FROM users WHERE id = $1", inviterID).Scan(&inviterName)
	var teamName string
	s.db.QueryRow("SELECT name FROM teams WHERE id = $1", teamID).Scan(&teamName)
	s.notifService.Create(inviteeID, "assigned", "Team Invitation", inviterName+" invited you to "+teamName, &inv.ID)

	return &inv, nil
}

// GetPending returns pending invitations for a user
func (s *InvitationService) GetPending(userID string) ([]models.Invitation, error) {
	rows, err := s.db.Query(
		`SELECT i.id, i.team_id, i.inviter_id, i.invitee_id, i.role, i.status, i.created_at, i.responded_at,
		        t.name, u.name, u.email
		 FROM team_invitations i
		 JOIN teams t ON t.id = i.team_id
		 JOIN users u ON u.id = i.inviter_id
		 WHERE i.invitee_id = $1 AND i.status = 'pending'
		 ORDER BY i.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []models.Invitation
	for rows.Next() {
		var inv models.Invitation
		var teamName, inviterName, inviterEmail string
		if err := rows.Scan(
			&inv.ID, &inv.TeamID, &inv.InviterID, &inv.InviteeID, &inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
			&teamName, &inviterName, &inviterEmail,
		); err != nil {
			return nil, err
		}
		inv.Team = &models.Team{Name: teamName}
		inv.Inviter = &models.User{Name: inviterName, Email: inviterEmail}
		invitations = append(invitations, inv)
	}
	if invitations == nil {
		invitations = []models.Invitation{}
	}
	return invitations, nil
}

// Accept accepts an invitation and adds user to team
func (s *InvitationService) Accept(invitationID, userID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get invitation
	var teamID, inviteeID, role string
	err = tx.QueryRow(
		"SELECT team_id, invitee_id, role FROM team_invitations WHERE id = $1 AND status = 'pending'",
		invitationID,
	).Scan(&teamID, &inviteeID, &role)
	if err != nil {
		return ErrInvitationNotFound
	}

	if inviteeID != userID {
		return ErrNotInvitee
	}

	// Update invitation status
	_, err = tx.Exec(
		"UPDATE team_invitations SET status = 'accepted', responded_at = NOW() WHERE id = $1",
		invitationID,
	)
	if err != nil {
		return err
	}

	// Add to team
	_, err = tx.Exec(
		"INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		teamID, userID, role,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Reject rejects an invitation
func (s *InvitationService) Reject(invitationID, userID string) error {
	result, err := s.db.Exec(
		"UPDATE team_invitations SET status = 'rejected', responded_at = NOW() WHERE id = $1 AND invitee_id = $2 AND status = 'pending'",
		invitationID, userID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrInvitationNotFound
	}
	return nil
}

// GetTeamInvitations returns all invitations for a team (for leaders)
func (s *InvitationService) GetTeamInvitations(teamID string) ([]models.Invitation, error) {
	rows, err := s.db.Query(
		`SELECT i.id, i.team_id, i.inviter_id, i.invitee_id, i.role, i.status, i.created_at, i.responded_at,
		        u.name, u.email
		 FROM team_invitations i
		 JOIN users u ON u.id = i.invitee_id
		 WHERE i.team_id = $1
		 ORDER BY i.created_at DESC`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []models.Invitation
	for rows.Next() {
		var inv models.Invitation
		var inviteeName, inviteeEmail string
		if err := rows.Scan(
			&inv.ID, &inv.TeamID, &inv.InviterID, &inv.InviteeID, &inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
			&inviteeName, &inviteeEmail,
		); err != nil {
			return nil, err
		}
		inv.Invitee = &models.User{Name: inviteeName, Email: inviteeEmail}
		invitations = append(invitations, inv)
	}
	if invitations == nil {
		invitations = []models.Invitation{}
	}
	return invitations, nil
}
