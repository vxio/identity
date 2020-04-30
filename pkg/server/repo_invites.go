package identityserver

import (
	"database/sql"
)

type InvitesRepository interface {
	list(tenantID string) ([]Invite, error)
	add(tenantID string, invite Invite, secretCode string) (Invite, error)
	delete(tenantID string, inviteID string) error
}

func NewInvitesRepository(db *sql.DB) InvitesRepository {
	return &sqlInvitesRepo{db: db}
}

type sqlInvitesRepo struct {
	db *sql.DB
}

func (r *sqlInvitesRepo) list(tenantID string) ([]Invite, error) {
	return nil, nil
}

func (r *sqlInvitesRepo) add(tenantID string, invite Invite, secretCode string) (Invite, error) {
	return invite, nil
}

func (r *sqlInvitesRepo) delete(tenantID string, inviteID string) error {
	return nil
}
