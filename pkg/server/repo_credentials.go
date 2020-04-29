package identityserver

import (
	"database/sql"
)

type CredentialRepository interface {
	list(tenantID string, identityID string) ([]Credential, error)
	add(tenantID string, identityID string, credentials Credential) (Credential, error)
}

func NewCredentialRepository(db *sql.DB) CredentialRepository {
	return &sqlCredsRepo{db: db}
}

type sqlCredsRepo struct {
	db *sql.DB
}

func (r *sqlCredsRepo) list(tenantID string, identityID string) ([]Credential, error) {
	return nil, nil
}

func (r *sqlCredsRepo) add(tenantID string, identityID string, credentials Credential) (Credential, error) {
	return credentials, nil
}
