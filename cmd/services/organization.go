package services

import (
	"database/sql"
	"errors"

	"kanbanmaster/cmd/models"
)

var (
	ErrOrgNotFound  = errors.New("organization not found")
	ErrNotOrgOwner  = errors.New("not organization owner")
)

type OrgService struct {
	db *sql.DB
}

func NewOrgService(db *sql.DB) *OrgService {
	return &OrgService{db: db}
}

type CreateOrgInput struct {
	Name string `json:"name"`
}

func (s *OrgService) Create(ownerID string, input CreateOrgInput) (*models.Organization, error) {
	var org models.Organization
	err := s.db.QueryRow(
		`INSERT INTO organizations (name, owner_id) VALUES ($1, $2)
		 RETURNING id, name, owner_id, created_at`,
		input.Name, ownerID,
	).Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (s *OrgService) List(userID string) ([]models.Organization, error) {
	rows, err := s.db.Query(
		`SELECT DISTINCT o.id, o.name, o.owner_id, o.created_at
		 FROM organizations o
		 LEFT JOIN teams t ON t.organization_id = o.id
		 LEFT JOIN team_members tm ON tm.team_id = t.id
		 WHERE o.owner_id = $1 OR tm.user_id = $1
		 ORDER BY o.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (s *OrgService) Get(id string) (*models.Organization, error) {
	var org models.Organization
	err := s.db.QueryRow(
		"SELECT id, name, owner_id, created_at FROM organizations WHERE id = $1",
		id,
	).Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (s *OrgService) Update(id, ownerID string, input CreateOrgInput) (*models.Organization, error) {
	var org models.Organization
	err := s.db.QueryRow(
		`UPDATE organizations SET name = $1 WHERE id = $2 AND owner_id = $3
		 RETURNING id, name, owner_id, created_at`,
		input.Name, id, ownerID,
	).Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotOrgOwner
		}
		return nil, err
	}
	return &org, nil
}

func (s *OrgService) Delete(id, ownerID string) error {
	result, err := s.db.Exec(
		"DELETE FROM organizations WHERE id = $1 AND owner_id = $2",
		id, ownerID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotOrgOwner
	}
	return nil
}
