package credentials

import (
	"database/sql"

	api "github.com/moov-io/identity/pkg/api"
)

type CredentialRepository interface {
	list(identityID string) ([]api.Credential, error)
	lookup(providerID string, subjectID string) (*api.Credential, error)
	get(credentialID string) (*api.Credential, error)
	add(credentials api.Credential) (*api.Credential, error)
	update(updated api.Credential) (*api.Credential, error)
}

func NewCredentialRepository(db *sql.DB) CredentialRepository {
	return &sqlCredsRepo{db: db}
}

type sqlCredsRepo struct {
	db *sql.DB
}

func (r *sqlCredsRepo) list(identityID string) ([]api.Credential, error) {
	qry := `
		SELECT credential_id, provider, subject_id, identity_id, creaton_on, last_used_on, disabled_on, disabled_by
		FROM credentials
		WHERE
			identity_id = ?
	`

	return r.queryScan(qry, identityID)
}

func (r *sqlCredsRepo) lookup(providerID string, subjectID string) (*api.Credential, error) {
	qry := `
		SELECT credential_id, provider, subject_id, identity_id, created_on, last_used_on, disabled_on, disabled_by
		FROM credentials
		WHERE provider_id = ? AND subject_id = ?
		LIMIT 1
	`

	results, err := r.queryScan(qry, providerID, subjectID)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, sql.ErrNoRows
	}

	return &results[0], nil
}

func (r *sqlCredsRepo) get(credentialID string) (*api.Credential, error) {
	qry := `
		SELECT credential_id, provider, subject_id, identity_id, created_on, last_used_on, disabled_on, disabled_by
		FROM 
		WHERE credential_id = ?
		LIMIT 1
	`

	results, err := r.queryScan(qry, credentialID)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, sql.ErrNoRows
	}

	return &results[0], nil
}

func (r *sqlCredsRepo) add(credentials api.Credential) (*api.Credential, error) {
	qry := `
		INSERT INTO credentials(
			credential_id, 
			provider, 
			subject_id, 
			identity_id, 
			created_on, 
			last_used_on, 
			disabled_on, 
			disabled_by)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(qry,
		credentials.CredentialID,
		credentials.Provider,
		credentials.SubjectID,
		credentials.IdentityID,
		credentials.CreatedOn,
		credentials.LastUsedOn,
		credentials.DisabledOn,
		credentials.DisabledBy)

	if err != nil {
		return nil, err
	}

	if cnt, err := res.RowsAffected(); cnt != 1 || err != nil {
		return nil, sql.ErrNoRows
	}

	return &credentials, nil
}

func (r *sqlCredsRepo) update(updated api.Credential) (*api.Credential, error) {

	qry := `
		UPDATE credentials
		SET
			last_used_on = ?
			disabled_on = ?
			disabled_by = ?
		WHERE
			credential_id = ? AND
			provider = ? AND
			subject_id = ? AND
			identity_id = ?
	`
	res, err := r.db.Exec(qry,
		updated.LastUsedOn,
		updated.DisabledOn,
		updated.DisabledBy,
		updated.CredentialID,
		updated.Provider,
		updated.SubjectID,
		updated.IdentityID)

	if err != nil {
		return nil, err
	}

	if cnt, err := res.RowsAffected(); cnt != 1 || err != nil {
		return nil, sql.ErrNoRows
	}

	return &updated, nil
}

func (r *sqlCredsRepo) queryScan(query string, args ...interface{}) ([]api.Credential, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := []api.Credential{}
	for rows.Next() {
		cred := api.Credential{}
		if err := rows.Scan(&cred.CredentialID, &cred.Provider, &cred.SubjectID, &cred.IdentityID, &cred.CreatedOn, &cred.LastUsedOn, &cred.DisabledOn, &cred.DisabledBy); err != nil {
			return nil, err
		}

		credentials = append(credentials, cred)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return credentials, nil
}
