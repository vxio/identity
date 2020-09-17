package identitiestestutils

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	fuzz "github.com/google/gofuzz"
	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/client"
)

func RandStringN(c fuzz.Continue, n int) string {
	s := make([]rune, n)
	for i := range s {
		c.Fuzz(&s[i])
	}
	return string(s)
}

func RandBirthDate(c fuzz.Continue) string {
	return time.Unix(c.Int63n(time.Now().Unix()), 0).In(time.UTC).Format(time.RFC3339)
}

func RandCountryCode(c fuzz.Continue) string {
	n := c.Intn(len(govalidator.ISO3166List))
	return govalidator.ISO3166List[n].Alpha2Code
}

func RandPostalCode(c fuzz.Continue) string {
	if c.RandBool() {
		return fmt.Sprintf("%d-%d", 10000+c.Intn(8999), 1000+c.Intn(899))
	} else {
		return fmt.Sprintf("%d", 10000+c.Intn(8999))
	}
}

func RandPhoneNumber(c fuzz.Continue) string {
	return fmt.Sprintf("%d-%d-%d", 100+c.Intn(89), 100+c.Intn(89), 1000+c.Intn(899))
}

func RandEmail(c fuzz.Continue) string {
	return fmt.Sprintf("test.%d@moov.io", 10000+c.Intn(8999))
}

func RandAddress(c fuzz.Continue) string {
	return fmt.Sprintf("%d %s St", 1000+c.Intn(899), c.RandString())
}

func RandNullable(c fuzz.Continue, value string) *string {
	if c.RandBool() {
		return nil
	} else {
		return &value
	}
}

func NewFuzzer() *fuzz.Fuzzer {

	return fuzz.New().NumElements(0, 5).Funcs(
		func(e *client.RegisterAddress, c fuzz.Continue) {
			e.Type = "primary"
			e.Address1 = RandAddress(c)
			e.Address2 = RandNullable(c, RandAddress(c))
			e.City = "cty" + c.RandString()
			e.State = RandCountryCode(c)
			e.PostalCode = RandPostalCode(c)
			e.Country = RandCountryCode(c)
		},

		func(e *client.UpdateAddress, c fuzz.Continue) {
			e.Type = "primary"
			e.Address1 = RandAddress(c)
			e.Address2 = RandNullable(c, RandAddress(c))
			e.City = "cty" + c.RandString()
			e.State = RandCountryCode(c)
			e.PostalCode = RandPostalCode(c)
			e.Country = RandCountryCode(c)
			e.Validated = c.RandBool()
		},

		func(e *client.Address, c fuzz.Continue) {
			e.Type = "primary"
			e.Address1 = RandAddress(c)
			e.Address2 = RandNullable(c, RandAddress(c))
			e.City = "cty" + c.RandString()
			e.State = RandCountryCode(c)
			e.PostalCode = RandPostalCode(c)
			e.Country = RandCountryCode(c)
		},

		func(e *client.RegisterPhone, c fuzz.Continue) {
			e.Type = "mobile"
			e.Number = RandPhoneNumber(c)
		},

		func(e *client.UpdatePhone, c fuzz.Continue) {
			e.Type = "mobile"
			e.Number = RandPhoneNumber(c)
			e.Validated = c.RandBool()
		},

		func(e *client.Phone, c fuzz.Continue) {
			e.Type = "mobile"
			e.Number = RandPhoneNumber(c)
		},

		func(e *client.Register, c fuzz.Continue) {
			e.CredentialID = uuid.New().String()

			e.InviteCode = c.RandString()
			e.FirstName = "fn" + c.RandString()
			e.MiddleName = c.RandString()
			e.LastName = "ln" + c.RandString()
			c.Fuzz(e.NickName)
			c.Fuzz(e.Suffix)
			e.BirthDate = RandNullable(c, RandBirthDate(c))
			e.Email = RandEmail(c)

			e.Phones = make([]client.RegisterPhone, c.Intn(3)+1)
			for i := range e.Phones {
				c.Fuzz(&e.Phones[i])
			}

			e.Addresses = make([]client.RegisterAddress, c.Intn(3)+1)
			for i := range e.Addresses {
				c.Fuzz(&e.Addresses[i])
			}
		},

		func(e *client.UpdateIdentity, c fuzz.Continue) {
			e.FirstName = "fn" + c.RandString()
			e.MiddleName = c.RandString()
			e.LastName = "ln" + c.RandString()
			c.Fuzz(e.NickName)
			c.Fuzz(e.Suffix)
			e.BirthDate = RandNullable(c, RandBirthDate(c))
			e.Status = RandStringN(c, 10)

			e.Phones = make([]client.UpdatePhone, c.Intn(3)+1)
			for i := range e.Phones {
				c.Fuzz(&e.Phones[i])
			}

			e.Addresses = make([]client.UpdateAddress, c.Intn(3)+1)
			for i := range e.Addresses {
				c.Fuzz(&e.Addresses[i])
			}
		},
	)
}
