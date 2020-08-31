package client

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (a Register) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CredentialID, validation.Required, is.UUID),
		validation.Field(&a.TenantID, is.UUID),
		validation.Field(&a.InviteCode, validation.Length(1, 60)),
		validation.Field(&a.FirstName, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.MiddleName, validation.Length(1, 255)),
		validation.Field(&a.LastName, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.Suffix, validation.Length(2, 20)),
		validation.Field(&a.BirthDate, validation.Date(time.RFC3339).Max(time.Now())),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.ImageUrl, is.URL),
		validation.Field(&a.Phones),
		validation.Field(&a.Addresses),
	)
}

func (a RegisterPhone) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Number, validation.Required, validation.Match(regexp.MustCompile("^[0-9-]{10,15}$"))),
		validation.Field(&a.Type, validation.Required, validation.In("home", "work", "mobile")),
	)
}

func (a RegisterAddress) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Type, validation.Required, validation.In("primary", "secondary")),
		validation.Field(&a.Address1, validation.Required, validation.Length(4, 255)),
		validation.Field(&a.Address2, validation.Length(4, 255)),
		validation.Field(&a.City, validation.Required, validation.Length(2, 255)),
		validation.Field(&a.State, validation.Required, validation.Length(2, 2)),
		validation.Field(&a.PostalCode, validation.Required, validation.Match(regexp.MustCompile("^[0-9-]{5,10}$"))),
		validation.Field(&a.Country, validation.Required, is.CountryCode2),
	)
}
