package identityserver

import (
	"database/sql"
	"errors"
	"fmt"
)

type InvitesRepository interface {
	list(tenantID string) ([]Invite, error)
	add(invite Invite, secretCode string) (*Invite, error)
	delete(tenantID string, inviteID string) error
}

func NewInvitesRepository(db *sql.DB) InvitesRepository {
	return &sqlInvitesRepo{db: db}
}

type sqlInvitesRepo struct {
	db *sql.DB
}

func (r *sqlInvitesRepo) list(tenantID string) ([]Invite, error) {
	qry := fmt.Sprintf(`
		SELECT %s FROM invites WHERE tenant_id = ?
	`, inviteSelect)

	return r.queryScan(qry, tenantID)
}

func (r *sqlInvitesRepo) add(invite Invite, secretCode string) (*Invite, error) {
	qry := `
		INSERT INTO invites(
			invite_id,
			tenant_id,
			email,
			invited_by,
			invited_on,
			expires_on,
			redeemed,
			secret_code
		) VALUES (?,?,?,?,?,?,?)`

	_, err := r.db.Exec(qry,
		invite.InviteID,
		invite.TenantID,
		invite.Email,
		invite.InvitedBy,
		invite.InvitedOn,
		invite.ExpiresOn,
		invite.Redeemed,
		secretCode)

	if err != nil {
		return nil, err
	}

	return &invite, nil
}

func (r *sqlInvitesRepo) delete(tenantID string, inviteID string) error {
	qry := `DELETE FROM invites WHERE tenant_id = ? AND invite_id = ?`

	res, err := r.db.Exec(qry, tenantID, inviteID)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("Invite not found to be deleted")
	}

	return nil
}

// Matches the order pulled in by the rows.Scan below in queryScanIdentity
var inviteSelect = `
	invites.identity_id,
	invites.tenant_id,
	invites.email,
	invites.invited_by,
	invites.invited_on,
	invites.expires_on,
	invites.redeemed
`

func (r *sqlInvitesRepo) queryScan(query string, args ...interface{}) ([]Invite, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Invite{}
	for rows.Next() {
		item := Invite{}
		if err := rows.Scan(&item.InviteID, &item.TenantID, &item.Email, &item.InvitedBy, &item.InvitedOn, &item.ExpiresOn, &item.Redeemed); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
