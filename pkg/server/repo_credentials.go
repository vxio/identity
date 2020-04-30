package identityserver

import (
	"database/sql"
)

type CredentialRepository interface {
	list(tenantID string, identityID string) ([]Credential, error)
	add(credentials Credential) (Credential, error)
	lookup(providerID string, subjectID string) (Credential, error)
	get(credentialID string) (Credential, error)
	update(updated Credential) (Credential, error)
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

func (r *sqlCredsRepo) add(credentials Credential) (Credential, error) {
	return credentials, nil
}

func (r *sqlCredsRepo) lookup(providerID string, subjectID string) (Credential, error) {
	// demo one for now.
	return Credential{
		Provider:  providerID,
		SubjectID: subjectID,
	}, nil
}

func (r *sqlCredsRepo) get(credentialID string) (Credential, error) {
	return Credential{
		CredentialID: credentialID,
	}, nil
}

func (r *sqlCredsRepo) update(updated Credential) (Credential, error) {
	return updated, nil
}
