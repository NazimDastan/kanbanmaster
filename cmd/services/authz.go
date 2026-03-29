package services

import "database/sql"

type AuthzService struct {
	db *sql.DB
}

func NewAuthzService(db *sql.DB) *AuthzService {
	return &AuthzService{db: db}
}

// UserCanAccessBoard checks if user is a member of the board's team
func (s *AuthzService) UserCanAccessBoard(userID, boardID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM boards b
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE b.id = $1 AND tm.user_id = $2
	`, boardID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessTask checks if user can access a task through board -> team chain
func (s *AuthzService) UserCanAccessTask(userID, taskID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		JOIN columns c ON c.id = t.column_id
		JOIN boards b ON b.id = c.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE t.id = $1 AND tm.user_id = $2
	`, taskID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessColumn checks if user can access a column through board -> team chain
func (s *AuthzService) UserCanAccessColumn(userID, columnID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM columns c
		JOIN boards b ON b.id = c.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE c.id = $1 AND tm.user_id = $2
	`, columnID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessSubtask checks via subtask -> task -> column -> board -> team
func (s *AuthzService) UserCanAccessSubtask(userID, subtaskID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM subtasks st
		JOIN tasks t ON t.id = st.task_id
		JOIN columns c ON c.id = t.column_id
		JOIN boards b ON b.id = c.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE st.id = $1 AND tm.user_id = $2
	`, subtaskID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessTeam checks if user is a member of the team
func (s *AuthzService) UserCanAccessTeam(userID, teamID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM team_members
		WHERE team_id = $1 AND user_id = $2
	`, teamID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessOrg checks if user is owner or member of any team in the org
func (s *AuthzService) UserCanAccessOrg(userID, orgID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM organizations o
		LEFT JOIN teams t ON t.organization_id = o.id
		LEFT JOIN team_members tm ON tm.team_id = t.id AND tm.user_id = $2
		WHERE o.id = $1 AND (o.owner_id = $2 OR tm.user_id IS NOT NULL)
	`, orgID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessLabel checks via label -> board -> team chain
func (s *AuthzService) UserCanAccessLabel(userID, labelID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM labels l
		JOIN boards b ON b.id = l.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE l.id = $1 AND tm.user_id = $2
	`, labelID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessComment checks via comment -> task -> column -> board -> team chain
func (s *AuthzService) UserCanAccessComment(userID, commentID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM comments cm
		JOIN tasks t ON t.id = cm.task_id
		JOIN columns c ON c.id = t.column_id
		JOIN boards b ON b.id = c.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE cm.id = $1 AND tm.user_id = $2
	`, commentID, userID).Scan(&count)
	return count > 0, err
}

// UserCanAccessAttachment checks via attachment -> task -> column -> board -> team chain
func (s *AuthzService) UserCanAccessAttachment(userID, attachmentID string) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM attachments a
		JOIN tasks t ON t.id = a.task_id
		JOIN columns c ON c.id = t.column_id
		JOIN boards b ON b.id = c.board_id
		JOIN team_members tm ON tm.team_id = b.team_id
		WHERE a.id = $1 AND tm.user_id = $2
	`, attachmentID, userID).Scan(&count)
	return count > 0, err
}
