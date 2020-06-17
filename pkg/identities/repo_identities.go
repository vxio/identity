package identities

import (
	"database/sql"
	"fmt"

	api "github.com/moov-io/identity/pkg/api"
)

// Repository - Used for interacting identities on the data store
type Repository interface {
	list(tenantID api.TenantID) ([]api.Identity, error)
	get(identityID string) (*api.Identity, error)
	update(updated api.Identity) (*api.Identity, error)
	add(identity api.Identity) (*api.Identity, error)
}

// NewIdentityRepository - Builds a new repository tied to the DB passed in.
func NewIdentityRepository(db *sql.DB) Repository {
	return &sqlIdentityRepo{db: db}
}

type sqlIdentityRepo struct {
	db *sql.DB
}

func (r *sqlIdentityRepo) list(tenantID api.TenantID) ([]api.Identity, error) {
	qry := fmt.Sprintf(`
		SELECT %s
		FROM identity
		WHERE identity.tenant_id = ?
	`, identitySelect)

	identities, err := r.queryScanIdentity(qry, tenantID.String())
	if err != nil {
		return nil, err
	}

	if len(identities) == 0 {
		return identities, nil
	}

	qry = fmt.Sprintf(`
		SELECT %s
		FROM identity_address
		INNER JOIN identity ON identity_address.identity_id = identity.identity_id
		WHERE identity.tenant_id = ?
	`, addressSelect)

	addresses, err := r.queryScanAddresses(qry, tenantID.String())
	if err != nil {
		return nil, err
	}

	qry = fmt.Sprintf(`
		SELECT %s
		FROM identity_phone
		INNER JOIN identity ON identity_phone.identity_id = identity.identity_id
		WHERE identity.tenant_id = ?
	`, phoneSelect)

	phones, err := r.queryScanPhone(qry, tenantID.String())
	if err != nil {
		return nil, err
	}

	for idx := range identities {
		i := &identities[idx]

		for _, a := range addresses {
			if a.IdentityID == i.IdentityID {
				i.Addresses = append(i.Addresses, a)
			}
		}

		for _, p := range phones {
			if p.IdentityID == i.IdentityID {
				i.Phones = append(i.Phones, p)
			}
		}
	}

	return identities, nil
}

func (r *sqlIdentityRepo) get(identityID string) (*api.Identity, error) {

	qry := fmt.Sprintf(`
		SELECT %s
		FROM identity
		WHERE identity.identity_id = ?
		LIMIT 1
	`, identitySelect)

	identities, err := r.queryScanIdentity(qry, identityID)
	if err != nil {
		return nil, err
	}

	if len(identities) != 1 {
		return nil, sql.ErrNoRows
	}

	qry = fmt.Sprintf(`
		SELECT %s
		FROM identity_address
		INNER JOIN identity ON identity_address.identity_id = identity.identity_id
		WHERE identity.identity_id = ?
	`, addressSelect)

	addresses, err := r.queryScanAddresses(qry, identityID)
	if err != nil {
		return nil, err
	}

	qry = fmt.Sprintf(`
		SELECT %s
		FROM identity_phone
		INNER JOIN identity ON identity_phone.identity_id = identity.identity_id
		WHERE identity.identity_id = ?
	`, phoneSelect)

	phones, err := r.queryScanPhone(qry, identityID)
	if err != nil {
		return nil, err
	}

	if len(identities) != 1 {
		return nil, sql.ErrNoRows
	}

	identities[0].Phones = phones
	identities[0].Addresses = addresses

	return &identities[0], nil
}

func (r *sqlIdentityRepo) update(updated api.Identity) (*api.Identity, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qry := `
		UPDATE identity
		SET
			first_name = ?,
			middle_name = ?,
			last_name = ?,
			nick_name = ?,
			suffix = ?,
			birth_date = ?,
			status = ?,
			disabled_on = ?,
			disabled_by = ?,
			last_updated_on = ?
		WHERE
			tenant_id = ? AND
			identity_id = ?
	`

	res, err := tx.Exec(qry,

		updated.FirstName,
		updated.MiddleName,
		updated.LastName,
		updated.NickName,
		updated.Suffix,
		updated.BirthDate,
		updated.Status,
		updated.DisabledOn,
		updated.DisabledBy,
		updated.LastUpdatedOn,

		updated.TenantID,
		updated.IdentityID)
	if err != nil {
		return nil, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, sql.ErrNoRows
	}

	if err := r.upsertAddresses(tx, &updated); err != nil {
		return nil, err
	}

	if err := r.upsertPhones(tx, &updated); err != nil {
		return nil, err
	}

	tx.Commit()

	return &updated, nil
}

func (r *sqlIdentityRepo) add(identity api.Identity) (*api.Identity, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qry := `
		INSERT INTO identity(
			identity_id, 
			tenant_id, 
			first_name, 
			middle_name, 
			last_name, 
			nick_name, 
			suffix, 
			birth_date, 
			status, 
			email, 
			email_verified,
			registered_on,
			invite_id,
			disabled_on,
			disabled_by,
			last_updated_on
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	res, err := tx.Exec(qry,
		identity.IdentityID,
		identity.TenantID,
		identity.FirstName,
		identity.MiddleName,
		identity.LastName,
		identity.NickName,
		identity.Suffix,
		identity.BirthDate,
		identity.Status,
		identity.Email,
		identity.EmailVerified,
		identity.RegisteredOn,
		identity.InviteID,
		identity.DisabledOn,
		identity.DisabledBy,
		identity.LastUpdatedOn)
	if err != nil {
		return nil, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, sql.ErrNoRows
	}

	if err := r.upsertAddresses(tx, &identity); err != nil {
		return nil, err
	}

	if err := r.upsertPhones(tx, &identity); err != nil {
		return nil, err
	}

	tx.Commit()

	return &identity, nil
}

// Matches the order pulled in by the rows.Scan below in queryScanIdentity
var identitySelect = `
	identity.identity_id, 
	identity.tenant_id,
	identity.first_name, 
	identity.middle_name, 
	identity.last_name, 
	identity.nick_name, 
	identity.suffix, 
	identity.birth_date, 
	identity.status, 
	identity.email, 
	identity.email_verified, 
	identity.registered_on, 
	identity.invite_id,
	identity.disabled_on, 
	identity.disabled_by,
	identity.last_updated_on
`

func (r *sqlIdentityRepo) queryScanIdentity(query string, args ...interface{}) ([]api.Identity, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []api.Identity{}
	for rows.Next() {
		item := api.Identity{}
		if err := rows.Scan(
			&item.IdentityID,
			&item.TenantID,
			&item.FirstName,
			&item.MiddleName,
			&item.LastName,
			&item.NickName,
			&item.Suffix,
			&item.BirthDate,
			&item.Status,
			&item.Email,
			&item.EmailVerified,
			&item.RegisteredOn,
			&item.InviteID,
			&item.DisabledOn,
			&item.DisabledBy,
			&item.LastUpdatedOn,
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

// Matches the order pulled in by the rows.Scan below in queryScanIdentity
var addressSelect = `
	identity_address.identity_id,
	identity_address.address_id,
	identity_address.type, 
	identity_address.address_1, 
	identity_address.address_2, 
	identity_address.city, 
	identity_address.postal_code,
	identity_address.state, 
	identity_address.country, 
	identity_address.validated
`

func (r *sqlIdentityRepo) queryScanAddresses(query string, args ...interface{}) ([]api.Address, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []api.Address{}
	for rows.Next() {
		item := api.Address{}
		if err := rows.Scan(
			&item.IdentityID,
			&item.AddressID,
			&item.Type,
			&item.Address1,
			&item.Address2,
			&item.City,
			&item.PostalCode,
			&item.State,
			&item.Country,
			&item.Validated,
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

// Matches the order pulled in by the rows.Scan below in queryScanIdentity
var phoneSelect = `
	identity_phone.identity_id,
	identity_phone.phone_id,
	identity_phone.type,
	identity_phone.number,
	identity_phone.validated
`

func (r *sqlIdentityRepo) queryScanPhone(query string, args ...interface{}) ([]api.Phone, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []api.Phone{}
	for rows.Next() {
		item := api.Phone{}
		if err := rows.Scan(
			&item.IdentityID,
			&item.PhoneID,
			&item.Type,
			&item.Number,
			&item.Validated,
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

func (r *sqlIdentityRepo) upsertAddresses(tx *sql.Tx, updated *api.Identity) error {

	updateQry := `
		UPDATE identity_address
		SET 
			type = ?,
			address_1 = ?,
			address_2 = ?,
			city = ?,
			postal_code = ?,
			state = ?,
			country = ?,
			validated = ?,
			last_updated_on = ?
		WHERE
			identity_id = ? AND
			address_id = ?
	`

	for _, a := range updated.Addresses {
		c, err := tx.Exec(updateQry,
			a.Type,
			a.Address1,
			a.Address2,
			a.City,
			a.PostalCode,
			a.State,
			a.Country,
			a.Validated,
			updated.LastUpdatedOn,

			updated.IdentityID,
			a.AddressID)
		if err != nil {
			return err
		}

		cnt, err := c.RowsAffected()
		if err != nil {
			return err
		}

		if cnt == 0 {

			insertQry := `
				INSERT INTO identity_address (
					identity_id, 
					address_id, 
					type, 
					address_1, 
					address_2, 
					city, 
					postal_code,
					state, 
					country, 
					validated, 
					last_updated_on
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)	
			`

			_, err := tx.Exec(insertQry,
				updated.IdentityID,
				a.AddressID,
				a.Type,
				a.Address1,
				a.Address2,
				a.City,
				a.PostalCode,
				a.State,
				a.Country,
				a.Validated,
				updated.LastUpdatedOn,
			)
			if err != nil {
				return err
			}
		}
	}

	// cleanout non-updated addresses
	if _, err := tx.Exec(`DELETE FROM identity_address WHERE identity_id = ? AND last_updated_on < ?`, updated.IdentityID, updated.LastUpdatedOn); err != nil {
		return err
	}

	return nil
}

func (r *sqlIdentityRepo) upsertPhones(tx *sql.Tx, updated *api.Identity) error {

	updateQry := `
		UPDATE identity_phone
		SET 
			type = ?,
			number = ?,
			validated = ?,
			last_updated_on = ?
		WHERE
			identity_id = ? AND
			phone_id = ?
	`

	for _, p := range updated.Phones {
		c, err := tx.Exec(updateQry,
			p.Type,
			p.Number,
			p.Validated,
			updated.LastUpdatedOn,

			updated.IdentityID,
			p.PhoneID)
		if err != nil {
			return err
		}

		cnt, err := c.RowsAffected()
		if err != nil {
			return err
		}

		if cnt == 0 {

			insertQry := `
				INSERT INTO identity_phone (
					identity_id, 
					phone_id, 
					type, 
					number, 
					validated, 
					last_updated_on
				) VALUES (?, ?, ?, ?, ?, ?)	
			`

			_, err := tx.Exec(insertQry,
				updated.IdentityID,
				p.PhoneID,
				p.Type,
				p.Number,
				p.Validated,
				updated.LastUpdatedOn,
			)
			if err != nil {
				return err
			}
		}
	}

	// cleanout non-updated phones
	if _, err := tx.Exec(`DELETE FROM identity_phone WHERE identity_id = ? AND last_updated_on < ?`, updated.IdentityID, updated.LastUpdatedOn); err != nil {
		return err
	}

	return nil
}
