package identityserver

import (
	"database/sql"
)

type IdentityRepository interface {
	list(tenantID string) ([]Identity, error)
	get(tenantID string, identityID string) (Identity, error)
	update(tenantID string, updated Identity) (Identity, error)
	add(identity Identity) (Identity, error)
}

func NewIdentityRepository(db *sql.DB) IdentityRepository {
	return &sqlIdentityRepo{db: db}
}

type sqlIdentityRepo struct {
	db *sql.DB
}

func (r *sqlIdentityRepo) list(tenantID string) ([]Identity, error) {
	return nil, nil
}

func (r *sqlIdentityRepo) get(tenantID string, identityID string) (Identity, error) {
	found := Identity{}
	return found, nil
}

func (r *sqlIdentityRepo) update(tenantID string, updated Identity) (Identity, error) {
	return updated, nil
}

func (r *sqlIdentityRepo) add(identity Identity) (Identity, error) {
	return identity, nil
}
