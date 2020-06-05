package credentials

import (
	"database/sql"
	"fmt"
	"time"

	api "github.com/moov-io/identity/pkg/api"
)

type CredentialRepository interface {
	list(identityID string, tenantID string) ([]api.Credential, error)
	lookup(providerID string, subjectID string, tenantID string) (*api.Credential, error)
	get(identityID string, credentialID string, tenantID string) (*api.Credential, error)
	add(credentials api.Credential) (*api.Credential, error)
	update(updated api.Credential) (*api.Credential, error)
	record(credentialID string, nonce string, ip string, at time.Time) error
}

func NewCredentialRepository(db *sql.DB) CredentialRepository {
	return &sqlCredsRepo{db: db}
}

type sqlCredsRepo struct {
	db *sql.DB
}

func (r *sqlCredsRepo) list(identityID string, tenantID string) ([]api.Credential, error) {
	qry := fmt.Sprintf(`
		SELECT %s
		FROM credentials
		WHERE identity_id = ? AND tenant_id = ?
	`, credentialSelect)

	return r.queryScan(qry, identityID, tenantID)
}

func (r *sqlCredsRepo) lookup(providerID string, subjectID string, tenantID string) (*api.Credential, error) {
	qry := fmt.Sprintf(`
		SELECT %s
		FROM credentials
		WHERE provider = ? AND subject_id = ? AND tenant_id = ?
		LIMIT 1
	`, credentialSelect)

	results, err := r.queryScan(qry, providerID, subjectID, tenantID)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, sql.ErrNoRows
	}

	return &results[0], nil
}

func (r *sqlCredsRepo) get(identityID string, credentialID string, tenantID string) (*api.Credential, error) {
	qry := fmt.Sprintf(`
		SELECT %s
		FROM credentials
		WHERE credential_id = ? AND identity_id = ? AND tenant_id = ?
		LIMIT 1
	`, credentialSelect)

	results, err := r.queryScan(qry, credentialID, identityID, tenantID)
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
			tenant_id,
			provider, 
			subject_id, 
			identity_id,
			created_on, 
			last_used_on, 
			disabled_on, 
			disabled_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(qry,
		credentials.CredentialID,
		credentials.TenantID,
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

func (r *sqlCredsRepo) record(credentialID string, nonce string, ip string, at time.Time) error {
	qry := `
		INSERT INTO credential_logins(
			credential_id,
			nonce,
			ip,
			created_on
		) VALUES (?, ?, ?, ?)
	`

	res, err := r.db.Exec(qry,
		credentialID,
		nonce,
		ip,
		at,
	)

	if err != nil {
		return err
	}

	if cnt, err := res.RowsAffected(); cnt != 1 || err != nil {
		return sql.ErrNoRows
	}

	return nil
}

func (r *sqlCredsRepo) update(updated api.Credential) (*api.Credential, error) {

	qry := `
		UPDATE credentials
		SET
			last_used_on = ?,
			disabled_on = ?,
			disabled_by = ?
		WHERE
			credential_id = ? AND
			provider = ? AND
			subject_id = ? AND
			identity_id = ? AND
			tenant_id = ?
	`
	_, err := r.db.Exec(qry,
		updated.LastUsedOn,
		updated.DisabledOn,
		updated.DisabledBy,

		updated.CredentialID,
		updated.Provider,
		updated.SubjectID,
		updated.IdentityID,
		updated.TenantID)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

var credentialSelect = `
	credential_id, 
	tenant_id, 
	provider, 
	subject_id, 
	identity_id, 
	created_on, 
	last_used_on, 
	disabled_on, 
	disabled_by
`

func (r *sqlCredsRepo) queryScan(query string, args ...interface{}) ([]api.Credential, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := []api.Credential{}
	for rows.Next() {
		cred := api.Credential{}
		if err := rows.Scan(&cred.CredentialID, &cred.TenantID, &cred.Provider, &cred.SubjectID, &cred.IdentityID, &cred.CreatedOn, &cred.LastUsedOn, &cred.DisabledOn, &cred.DisabledBy); err != nil {
			return nil, err
		}

		credentials = append(credentials, cred)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return credentials, nil
}
