package invites

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/zerotrust"
)

func TestGetById(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository InvitesRepository) {
		invite, _ := AddTestingInvite(t, repository)

		tenantID := zerotrust.TenantID(uuid.MustParse(invite.TenantID))

		found, err := repository.get(tenantID, invite.InviteID)
		if err != nil {
			t.Error(err)
		}

		if *found != invite {
			t.Error("Found by ID doesn't match Invite", cmp.Diff(*found, invite))
		}

		badTenantID := zerotrust.TenantID(uuid.New())
		found, err = repository.get(badTenantID, invite.InviteID)
		if err != sql.ErrNoRows {
			t.Error(err)
		}
	})
}

func TestGetByCode(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository InvitesRepository) {
		invite, code := AddTestingInvite(t, repository)

		found, err := repository.getByCode(code)
		if err != nil {
			t.Error(err)
		}

		if *found != invite {
			t.Error("Found by code doesn't match invite", found, invite)
		}

		_, err = repository.getByCode("doesnotexist")
		if err != sql.ErrNoRows {
			t.Error(err)
		}
	})
}

func TestList(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository InvitesRepository) {
		invite, _ := AddTestingInvite(t, repository)
		tenantID1 := zerotrust.TenantID(uuid.MustParse(invite.TenantID))

		// Add noise and other invites on other tenants
		AddTestingInvite(t, repository)
		AddTestingInvite(t, repository)
		AddTestingInvite(t, repository)

		found, err := repository.list(tenantID1)
		if err != nil {
			t.Error(err)
		}

		if len(found) != 1 {
			t.Error("Found more than one invite on a tenant with only 1")
		}

		if found[0] != invite {
			t.Error("Invite from first tenant doesn't match list of first tenant")
		}

		badTenantID := zerotrust.TenantID(uuid.New())
		found, err = repository.list(badTenantID)
		if err != nil {
			t.Error(err)
		}

		if len(found) != 0 {
			t.Error("Returned rows on a bad tenant")
		}
	})
}

func TestUpdate(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository InvitesRepository) {
		invite, _ := AddTestingInvite(t, repository)
		tenantID := zerotrust.TenantID(uuid.MustParse(invite.TenantID))

		updated := invite
		redeemedOn := time.Now().In(time.UTC).Round(time.Second)
		disabledBy := uuid.New().String()
		updated.RedeemedOn = &redeemedOn
		updated.DisabledBy = &disabledBy
		updated.DisabledOn = &redeemedOn

		err := repository.update(updated)
		if err != nil {
			t.Error(err)
		}

		found, err := repository.get(tenantID, updated.InviteID)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(*found, updated) {
			t.Errorf("The updated found doesn't match the update\n%s", cmp.Diff(*found, updated))
		}

		badUpdate := updated
		badUpdate.TenantID = uuid.New().String()

		err = repository.update(badUpdate)
		if err != sql.ErrNoRows {
			t.Error(err)
		}
	})
}

func NewTestRepository(t *testing.T) InvitesRepository {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	err = database.RunMigrations(db, database.DatabaseConfig{
		SqlLite: &database.SqlLiteConfig{
			Path: ":memory:",
		},
	})

	if err != nil {
		panic(err)
	}

	repo := NewInvitesRepository(db)

	return repo
}

func ForEachDatabase(t *testing.T, run func(t *testing.T, repository InvitesRepository)) {
	cases := map[string]database.DatabaseConfig{
		"sqlite": database.DatabaseConfig{
			DatabaseName: "sqlite",
			SqlLite: &database.SqlLiteConfig{
				Path: ":memory:",
			},
		},
		"mysql": database.DatabaseConfig{
			DatabaseName: "identity",
			MySql: &database.MySqlConfig{
				Address:  "tcp(localhost:4306)",
				User:     "identity",
				Password: "identity",
			},
		},
	}

	for k, tc := range cases {
		t.Run(k, func(t *testing.T) {
			db, err := database.New(context.Background(), log.NewNopLogger(), tc)

			//db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				panic(err)
			}

			t.Cleanup(func() {
				db.Close()
			})

			err = database.RunMigrations(db, tc)

			if err != nil {
				panic(err)
			}

			repo := NewInvitesRepository(db)

			run(t, repo)
		})
	}
}

func AddTestingInvite(t *testing.T, repository InvitesRepository) (api.Invite, string) {
	i := RandomInvite()
	code, err := generateInviteCode()
	if err != nil {
		t.Error(err)
	}

	added, err := repository.add(i, *code)
	if err != nil {
		t.Error(err)
	}

	return *added, *code
}

func RandomInvite() api.Invite {
	return api.Invite{
		InviteID:   uuid.New().String(),
		TenantID:   uuid.New().String(),
		Email:      "someuser@domain.com",
		InvitedBy:  uuid.New().String(),
		InvitedOn:  time.Now().In(time.UTC).Round(time.Second),
		RedeemedOn: nil,
		ExpiresOn:  time.Now().Add(time.Hour).In(time.UTC).Round(time.Second),
		DisabledOn: nil,
		DisabledBy: nil,
	}
}
