package identityserver

import (
	"database/sql"
	"identityserver"
)

type IdentityRepository interface {
	list(tenantID string) ([]identityserver.Identity, error)
	get(tenantID string, identityID string) (identityserver.Identity, error)
	update(tenantID string, updated identityserver.Identity) (identityserver.Identity, error)
}

func NewIdentityRepository(db *sql.DB) IdentityRepository {
	return &sqlIdentityRepo{db: db}
}

type sqlIdentityRepo struct {
	db *sql.DB
}

func (r *sqlIdentityRepo) list(tenantID string) ([]identityserver.Identity, error) {
	return nil, nil
}

func (r *sqlIdentityRepo) get(tenantID string, identityID string) (identityserver.Identity, error) {
	found := identityserver.Identity()
	return found, nil
}

func (r *sqlIdentityRepo) update(tenantID string, updated identityserver.Identity) (identityserver.Identity, error) {
	return updated
}
