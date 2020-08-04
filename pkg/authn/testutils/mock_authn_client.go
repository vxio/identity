package authntestutils

import (
	"github.com/google/uuid"
	authnClient "github.com/moov-io/authn/pkg/client"
	authnclient "github.com/moov-io/identity/pkg/authn/client"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

type mockAuthnClient struct {
	tenant authnClient.Tenant
}

func NewMockAuthnClient() authnclient.AuthnClient {
	tenant := authnClient.Tenant{
		TenantID: uuid.New().String(),
		Name:     "My Tenant",
		Alias:    "my-tenant",
		Website:  "https://example.com",
	}

	return &mockAuthnClient{tenant: tenant}
}

func (mac *mockAuthnClient) GetTenant(claims tmw.TumblerClaims, tenantID string) (*authnClient.Tenant, error) {
	return &mac.tenant, nil
}
