package identityserver

import (
	"database/sql"
	"identityserver"
)

type CredentialRepository interface {
	list(tenantID string, identityID string) ([]identityserver.Credential, error)
	add(tenantID string, identityID string, credentials identityserver.Credential) (identityserver.Credential, error)
}

func NewCredentialRepository(db *sql.DB) CredentialRepository {
	return &sqlCredsRepo{db: db}
}

type sqlCredsRepo struct {
	db *sql.DB
}

func (r *sqlCredsRepo) list(tenantID string, identityID string) ([]identityserver.Credential, error) {
	return nil, nil
}

func (r *sqlCredsRepo) add(tenantID string, identityID string, credentials identityserver.Credential) (identityserver.Credential, error) {
	return credentials, nil
}
