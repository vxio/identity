package invites

import (
	"database/sql"
	"fmt"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/zerotrust"
)

// Repository allows for interacting with the invites data store.
type Repository interface {
	list(tenantID zerotrust.TenantID) ([]api.Invite, error)
	get(tenantID zerotrust.TenantID, inviteID string) (*api.Invite, error)
	getByCode(code string) (*api.Invite, error)
	add(invite api.Invite, secretCode string) (*api.Invite, error)
	update(updated api.Invite) error
}

// NewInvitesRepository instantiates a new InvitesRepository
func NewInvitesRepository(db *sql.DB) Repository {
	return &sqlInvitesRepo{db: db}
}

type sqlInvitesRepo struct {
	db *sql.DB
}

func (r *sqlInvitesRepo) list(tenantID zerotrust.TenantID) ([]api.Invite, error) {
	qry := fmt.Sprintf(`
		SELECT %s 
		FROM invites 
		WHERE 
			tenant_id = ?
		ORDER BY invites.invited_on DESC
	`, inviteSelect)

	return r.queryScan(qry, tenantID.String())
}

func (r *sqlInvitesRepo) get(tenantID zerotrust.TenantID, inviteID string) (*api.Invite, error) {
	qry := fmt.Sprintf(`
		SELECT %s
		FROM invites
		WHERE tenant_id = ? AND invite_id = ?
		LIMIT 1
	`, inviteSelect)

	res, err := r.queryScan(qry, tenantID.String(), inviteID)
	if err != nil {
		return nil, err
	}

	if len(res) != 1 {
		return nil, sql.ErrNoRows
	}

	return &res[0], nil
}

func (r *sqlInvitesRepo) getByCode(code string) (*api.Invite, error) {
	qry := fmt.Sprintf(`
		SELECT %s
		FROM invites
		WHERE secret_code = ?
		LIMIT 1
	`, inviteSelect)

	res, err := r.queryScan(qry, code)
	if err != nil {
		return nil, err
	}

	if len(res) != 1 {
		return nil, sql.ErrNoRows
	}

	return &res[0], nil
}

func (r *sqlInvitesRepo) add(invite api.Invite, secretCode string) (*api.Invite, error) {
	qry := `
		INSERT INTO invites(
			invite_id,
			tenant_id,
			email,
			invited_by,
			invited_on,
			redeemed_on,
			expires_on,
			secret_code
		) VALUES (?,?,?,?,?,?,?,?)`

	_, err := r.db.Exec(qry,
		invite.InviteID,
		invite.TenantID,
		invite.Email,
		invite.InvitedBy,
		invite.InvitedOn,
		invite.RedeemedOn,
		invite.ExpiresOn,
		secretCode)

	if err != nil {
		return nil, err
	}

	return &invite, nil
}

func (r *sqlInvitesRepo) update(updated api.Invite) error {
	qry := `
		UPDATE invites 
		SET 
			redeemed_on = ?,
			disabled_by = ?,
			disabled_on = ?
		WHERE
			tenant_id = ? AND 
			invite_id = ?
	`

	res, err := r.db.Exec(qry,
		updated.RedeemedOn,
		updated.DisabledBy,
		updated.DisabledOn,
		updated.TenantID,
		updated.InviteID)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Matches the order pulled in by the rows.Scan below in queryScanIdentity
var inviteSelect = `
	invites.invite_id,
	invites.tenant_id,
	invites.email,
	invites.invited_by,
	invites.invited_on,
	invites.redeemed_on,
	invites.expires_on,
	invites.disabled_on,
	invites.disabled_by
`

func (r *sqlInvitesRepo) queryScan(query string, args ...interface{}) ([]api.Invite, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []api.Invite{}
	for rows.Next() {
		item := api.Invite{}
		if err := rows.Scan(
			&item.InviteID,
			&item.TenantID,
			&item.Email,
			&item.InvitedBy,
			&item.InvitedOn,
			&item.RedeemedOn,
			&item.ExpiresOn,
			&item.DisabledOn,
			&item.DisabledBy,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
