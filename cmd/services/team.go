package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var (
	ErrTeamNotFound    = errors.New("team not found")
	ErrNotTeamLeader   = errors.New("insufficient team permissions")
	ErrAlreadyMember   = errors.New("user is already a team member")
	ErrMemberNotFound  = errors.New("team member not found")
)

type TeamService struct {
	db *sql.DB
}

func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{db: db}
}

type CreateTeamInput struct {
	Name           string `json:"name"`
	OrganizationID string `json:"organizationId"`
}

func (s *TeamService) Create(userID string, input CreateTeamInput) (*models.Team, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var team models.Team
	err = tx.QueryRow(
		`INSERT INTO teams (name, organization_id) VALUES ($1, $2)
		 RETURNING id, name, organization_id, created_at`,
		input.Name, input.OrganizationID,
	).Scan(&team.ID, &team.Name, &team.OrganizationID, &team.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Creator becomes leader
	_, err = tx.Exec(
		`INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, 'leader')`,
		team.ID, userID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *TeamService) List(userID string) ([]models.Team, error) {
	rows, err := s.db.Query(
		`SELECT t.id, t.name, t.organization_id, t.created_at
		 FROM teams t
		 JOIN team_members tm ON tm.team_id = t.id
		 WHERE tm.user_id = $1
		 ORDER BY t.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.OrganizationID, &team.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (s *TeamService) GetWithMembers(teamID string) (*models.Team, error) {
	var team models.Team
	err := s.db.QueryRow(
		"SELECT id, name, organization_id, created_at FROM teams WHERE id = $1",
		teamID,
	).Scan(&team.ID, &team.Name, &team.OrganizationID, &team.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	rows, err := s.db.Query(
		`SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.joined_at,
		        u.id, u.email, u.name, u.avatar_url, u.created_at, u.updated_at
		 FROM team_members tm
		 JOIN users u ON u.id = tm.user_id
		 WHERE tm.team_id = $1
		 ORDER BY tm.joined_at`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.TeamMember
		var u models.User
		err := rows.Scan(
			&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.JoinedAt,
			&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		m.User = &u
		team.Members = append(team.Members, m)
	}

	return &team, nil
}

func (s *TeamService) Update(teamID, userID string, name string) (*models.Team, error) {
	if err := s.checkPermission(teamID, userID); err != nil {
		return nil, err
	}

	var team models.Team
	err := s.db.QueryRow(
		`UPDATE teams SET name = $1 WHERE id = $2
		 RETURNING id, name, organization_id, created_at`,
		name, teamID,
	).Scan(&team.ID, &team.Name, &team.OrganizationID, &team.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *TeamService) Delete(teamID, userID string) error {
	if err := s.checkPermission(teamID, userID); err != nil {
		return err
	}
	_, err := s.db.Exec("DELETE FROM teams WHERE id = $1", teamID)
	return err
}

func (s *TeamService) InviteMember(teamID, userID, inviteeEmail, role string) (*models.TeamMember, error) {
	if err := s.checkPermission(teamID, userID); err != nil {
		return nil, err
	}

	// Find user by email
	var inviteeID string
	err := s.db.QueryRow("SELECT id FROM users WHERE email = $1", inviteeEmail).Scan(&inviteeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check not already member
	var exists bool
	err = s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)",
		teamID, inviteeID,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAlreadyMember
	}

	if role == "" {
		role = "member"
	}

	var member models.TeamMember
	err = s.db.QueryRow(
		`INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, $3)
		 RETURNING id, team_id, user_id, role, joined_at`,
		teamID, inviteeID, role,
	).Scan(&member.ID, &member.TeamID, &member.UserID, &member.Role, &member.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *TeamService) RemoveMember(teamID, userID, memberUserID string) error {
	if err := s.checkPermission(teamID, userID); err != nil {
		return err
	}
	result, err := s.db.Exec(
		"DELETE FROM team_members WHERE team_id = $1 AND user_id = $2",
		teamID, memberUserID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrMemberNotFound
	}
	return nil
}

func (s *TeamService) UpdateMemberRole(teamID, userID, memberUserID, newRole string) error {
	if err := s.checkPermission(teamID, userID); err != nil {
		return err
	}
	result, err := s.db.Exec(
		"UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3",
		newRole, teamID, memberUserID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrMemberNotFound
	}
	return nil
}

func (s *TeamService) GetUserRole(teamID, userID string) (string, error) {
	var role string
	err := s.db.QueryRow(
		"SELECT role FROM team_members WHERE team_id = $1 AND user_id = $2",
		teamID, userID,
	).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrMemberNotFound
		}
		return "", err
	}
	return role, nil
}

func (s *TeamService) checkPermission(teamID, userID string) error {
	role, err := s.GetUserRole(teamID, userID)
	if err != nil {
		return ErrNotTeamLeader
	}
	if role != "owner" && role != "leader" {
		return ErrNotTeamLeader
	}
	return nil
}
