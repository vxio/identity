package identities_test

import (
	"context"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/client"
)

func Test_Register(t *testing.T) {
	a, s, f := Setup(t)

	invite := s.RandomInvite()
	r := client.Register{}
	f.Fuzz(&r)
	r.TenantID = invite.TenantID
	imageUrl := "http://kittens.com/img123.jpg"
	r.ImageUrl = &imageUrl

	r.Phones = make([]client.RegisterPhone, 1)
	f.Fuzz(&r.Phones[0])

	r.Addresses = make([]client.RegisterAddress, 1)
	f.Fuzz(&r.Addresses[0])

	i, err := s.service.Register(r, &invite)
	a.Nil(err)

	a.Equal(invite.TenantID, i.TenantID)
	a.Equal(s.session.TenantID.String(), i.TenantID)

	a.Equal(r.FirstName, i.FirstName)
	a.Equal(r.MiddleName, i.MiddleName)
	a.Equal(r.LastName, i.LastName)
	a.Equal(r.NickName, i.NickName)
	a.Equal(r.Suffix, i.Suffix)
	a.Equal(r.BirthDate, i.BirthDate)
	a.Equal(r.Email, i.Email)
	a.Equal(r.ImageUrl, i.ImageUrl)

	a.Len(i.Phones, 1)
	for x, rp := range r.Phones {
		ip := i.Phones[x]
		a.Equal(rp.Type, ip.Type)
		a.Equal(rp.Number, ip.Number)
	}

	a.Len(i.Addresses, 1)
	for x, ra := range r.Addresses {
		ia := i.Addresses[x]
		a.Equal(ra.Address1, ia.Address1)
		a.Equal(ra.Address2, ia.Address2)
		a.Equal(ra.City, ia.City)
		a.Equal(ra.Country, ia.Country)
		a.Equal(ra.PostalCode, ia.PostalCode)
		a.Equal(ra.State, ia.State)
		a.Equal(ra.Type, ia.Type)
	}

	// Fail on second register
	_, err = s.service.Register(r, &invite)
	a.NotNil(err)
}

func Test_GetAPI(t *testing.T) {
	a, s, f := Setup(t)

	identity := RegisterIdentity(s, f)

	found, resp, err := s.api.IdentitiesApi.GetIdentity(context.Background(), identity.IdentityID)

	a.Nil(err)
	a.Equal(200, resp.StatusCode)
	a.Equal(identity, found)
	a.Equal(identity.ImageUrl, found.ImageUrl)
}

func Test_GetAPI_NotFound(t *testing.T) {
	a, s, _ := Setup(t)

	_, resp, _ := s.api.IdentitiesApi.GetIdentity(context.Background(), uuid.New().String())
	a.Equal(404, resp.StatusCode)
}

func Test_ListAPI(t *testing.T) {
	a, s, f := Setup(t)

	identity1 := RegisterIdentity(s, f)
	identity2 := RegisterIdentity(s, f)
	identity3 := RegisterIdentity(s, f)

	found, resp, err := s.api.IdentitiesApi.ListIdentities(context.Background(), nil)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	a.Len(found, 3)
	a.Contains(found, identity1, identity2, identity3)
}

func Test_ListAPI_Empty(t *testing.T) {
	a, s, _ := Setup(t)

	found, resp, err := s.api.IdentitiesApi.ListIdentities(context.Background(), nil)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	a.Len(found, 0)
}

func Test_UpdateAPI(t *testing.T) {
	a, s, f := Setup(t)

	identity := RegisterIdentity(s, f)

	s.time.Add(time.Millisecond)

	updates := client.UpdateIdentity{}
	f.Fuzz(&updates)

	updated, resp, err := s.api.IdentitiesApi.UpdateIdentity(
		context.Background(),
		identity.IdentityID,
		updates,
	)

	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// These shouldn't change.
	a.Equal(identity.IdentityID, updated.IdentityID)
	a.Equal(identity.TenantID, updated.TenantID)
	a.Equal(identity.Email, updated.Email)
	a.Equal(identity.RegisteredOn, updated.RegisteredOn)
	a.Equal(identity.LastLogin, updated.LastLogin)
	a.Equal(identity.DisabledOn, updated.DisabledOn)
	a.Equal(identity.DisabledBy, updated.DisabledBy)
	a.Equal(identity.InviteID, updated.InviteID)

	// These change
	a.Equal(updated.LastUpdatedOn, s.time.Now())

	a.Equal(updates.FirstName, updated.FirstName)
	a.Equal(updates.MiddleName, updated.MiddleName)
	a.Equal(updates.LastName, updated.LastName)
	a.Equal(updates.NickName, updated.NickName)
	a.Equal(updates.Suffix, updated.Suffix)
	a.Equal(updates.BirthDate, updated.BirthDate)
	a.Equal(updates.Status, updated.Status)

	a.Len(updated.Phones, len(updates.Phones))
	for idx := range updated.Phones {
		exp := updates.Phones[idx]
		cur := updated.Phones[idx]

		a.Equal(identity.IdentityID, cur.IdentityID)
		a.Equal(exp.Number, cur.Number)
		a.Equal(exp.Validated, cur.Validated)
		a.Equal(exp.Type, cur.Type)
	}

	a.Len(updated.Addresses, len(updates.Addresses))
	for idx := range updated.Addresses {
		exp := updates.Addresses[idx]
		cur := updated.Addresses[idx]

		a.Equal(identity.IdentityID, cur.IdentityID)
		a.Equal(exp.Type, cur.Type)
		a.Equal(exp.Address1, cur.Address1)
		a.Equal(exp.Address2, cur.Address2)
		a.Equal(exp.City, cur.City)
		a.Equal(exp.State, cur.State)
		a.Equal(exp.PostalCode, cur.PostalCode)
		a.Equal(exp.Country, cur.Country)
		a.Equal(exp.Validated, cur.Validated)
	}

	found, resp, err := s.api.IdentitiesApi.GetIdentity(context.Background(), identity.IdentityID)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	a.Len(found.Phones, len(updates.Phones))
	a.Len(found.Addresses, len(updates.Addresses))
	a.Equal(updated, found)
}

func Test_UpdateAPI_NotFound(t *testing.T) {
	a, s, f := Setup(t)

	updates := client.UpdateIdentity{}
	f.Fuzz(&updates)

	_, resp, _ := s.api.IdentitiesApi.UpdateIdentity(
		context.Background(),
		uuid.New().String(),
		updates,
	)

	a.Equal(404, resp.StatusCode)
}

func Test_DisableAPI(t *testing.T) {
	a, s, f := Setup(t)

	identity := RegisterIdentity(s, f)

	resp, err := s.api.IdentitiesApi.DisableIdentity(context.Background(), identity.IdentityID)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	disabled, resp, err := s.api.IdentitiesApi.GetIdentity(context.Background(), identity.IdentityID)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	a.Equal(s.time.Now(), *disabled.DisabledOn)
	a.Equal(s.time.Now(), disabled.LastUpdatedOn)
	a.Equal(s.session.IdentityID.String(), *disabled.DisabledBy)
}

func Test_DisableAPI_NotFound(t *testing.T) {
	a, s, _ := Setup(t)

	resp, _ := s.api.IdentitiesApi.DisableIdentity(context.Background(), uuid.New().String())
	a.Equal(404, resp.StatusCode)
}

func RegisterIdentity(s Scope, f *fuzz.Fuzzer) client.Identity {
	invite := s.RandomInvite()

	register := client.Register{}
	f.Fuzz(&register)

	identity, err := s.service.Register(register, &invite)
	if err != nil {
		panic(err)
	}

	return *identity
}
