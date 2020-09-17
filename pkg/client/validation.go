package client

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// OpenID Connect allows for a birth date to be just a year.
// So we nil this out so it doesn't get created as a real birth date
func fixBirthDate(dob *string) *string {
	if dob != nil && len(*dob) == 4 {
		return nil
	} else {
		return dob
	}
}

func validateBirthDate(value interface{}) error {
	r0 := validation.Date(time.RFC3339).Max(time.Now())
	if r0.Validate(value) == nil {
		return nil
	}

	r1 := validation.Date("2006-01-02").Max(time.Now())
	if r1.Validate(value) == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("invalid birthdate specified: %+v", value))
}

func (a *Register) Validate() error {
	a.BirthDate = fixBirthDate(a.BirthDate)
	return validation.ValidateStruct(a,
		validation.Field(&a.CredentialID, validation.Required, is.UUID),
		validation.Field(&a.TenantID, is.UUID),
		validation.Field(&a.InviteCode, validation.Length(1, 60)),
		validation.Field(&a.FirstName, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.MiddleName, validation.Length(1, 255)),
		validation.Field(&a.LastName, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.Suffix, validation.Length(2, 20)),
		validation.Field(&a.BirthDate, validation.By(validateBirthDate)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.ImageUrl, is.URL),
		validation.Field(&a.Phones),
		validation.Field(&a.Addresses),
	)
}

func (a *RegisterPhone) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Number, validation.Match(regexp.MustCompile("^[0-9-]{10,15}$"))),
		validation.Field(&a.Type, validation.In("home", "work", "mobile")),
	)
}

func (a *RegisterAddress) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Type, validation.In("primary", "secondary")),
		validation.Field(&a.Address1, validation.Length(4, 255)),
		validation.Field(&a.Address2, validation.Length(4, 255)),
		validation.Field(&a.City, validation.Length(2, 255)),
		validation.Field(&a.State, validation.Length(2, 2)),
		validation.Field(&a.PostalCode, validation.Match(regexp.MustCompile("^[0-9-]{5,10}$"))),
		validation.Field(&a.Country, is.CountryCode2),
	)
}

func (a *UpdateIdentity) Validate() error {
	a.BirthDate = fixBirthDate(a.BirthDate)
	return validation.ValidateStruct(a,
		validation.Field(&a.FirstName, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.MiddleName, validation.Length(1, 255)),
		validation.Field(&a.LastName, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.Suffix, validation.Length(2, 20)),
		validation.Field(&a.BirthDate, validation.By(validateBirthDate)),
		validation.Field(&a.Phones),
		validation.Field(&a.Addresses),
	)
}

func (a *UpdatePhone) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Number, validation.Required, validation.Match(regexp.MustCompile("^[0-9-]{10,15}$"))),
		validation.Field(&a.Type, validation.Required, validation.In("home", "work", "mobile")),
	)
}

func (a *UpdateAddress) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Type, validation.Required, validation.In("primary", "secondary")),
		validation.Field(&a.Address1, validation.Required, validation.Length(4, 255)),
		validation.Field(&a.Address2, validation.Length(4, 255)),
		validation.Field(&a.City, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.State, validation.Required, validation.Length(2, 2)),
		validation.Field(&a.PostalCode, validation.Required, validation.Match(regexp.MustCompile("^[0-9-]{5,10}$"))),
		validation.Field(&a.Country, validation.Required, is.CountryCode2),
	)
}
