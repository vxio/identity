package invites

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moov-io/base/docker"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/database"
	log "github.com/moov-io/identity/pkg/logging"
)

func TestGetById(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository Repository) {
		invite, _ := AddTestingInvite(t, repository)

		tenantID := api.TenantID(uuid.MustParse(invite.TenantID))

		found, err := repository.get(tenantID, invite.InviteID)
		if err != nil {
			t.Error(err)
		}

		if *found != invite {
			t.Error("found by ID doesn't match Invite", cmp.Diff(*found, invite))
		}

		badTenantID := api.TenantID(uuid.New())
		_, err = repository.get(badTenantID, invite.InviteID)
		if err != sql.ErrNoRows {
			t.Error(err)
		}
	})
}

func TestGetByCode(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository Repository) {
		invite, code := AddTestingInvite(t, repository)

		found, err := repository.getByCode(code)
		if err != nil {
			t.Error(err)
		}

		if *found != invite {
			t.Error("found by code doesn't match invite", found, invite)
		}

		_, err = repository.getByCode("doesnotexist")
		if err != sql.ErrNoRows {
			t.Error(err)
		}
	})
}

func TestList(t *testing.T) {
	ForEachDatabase(t, func(t *testing.T, repository Repository) {
		invite, _ := AddTestingInvite(t, repository)
		tenantID1 := api.TenantID(uuid.MustParse(invite.TenantID))

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

		badTenantID := api.TenantID(uuid.New())
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
	ForEachDatabase(t, func(t *testing.T, repository Repository) {
		invite, _ := AddTestingInvite(t, repository)
		tenantID := api.TenantID(uuid.MustParse(invite.TenantID))

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

func ForEachDatabase(t *testing.T, run func(t *testing.T, repository Repository)) {
	cases := map[string]database.DatabaseConfig{
		"sqlite": database.InMemorySqliteConfig,
		"mysql": {
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
			if !docker.Enabled() && tc != database.InMemorySqliteConfig {
				t.SkipNow()
			}

			db := LoadDatabase(t, tc)
			repo := NewInvitesRepository(db)
			run(t, repo)
		})
	}
}

func LoadDatabase(t *testing.T, config database.DatabaseConfig) *sql.DB {
	l := log.NewNopLogger()
	db, err := database.New(context.Background(), l, config)
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	err = database.RunMigrations(l, db, config)
	if err != nil {
		panic(err)
	}

	return db
}

func NewInMemoryInvitesRepository(t *testing.T) Repository {
	db, close, err := database.NewAndMigrate(database.InMemorySqliteConfig, log.NewNopLogger(), context.Background())
	t.Cleanup(close)
	if err != nil {
		t.Error(err)
	}

	repo := NewInvitesRepository(db)

	return repo
}

func AddTestingInvite(t *testing.T, repository Repository) (api.Invite, string) {
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
