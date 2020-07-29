package authn

import (
	"net/http"
	"net/url"

	authnClient "github.com/moov-io/authn/pkg/client"

	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

type AuthnClient interface {
	GetTenant(claims tmw.TumblerClaims, tenantID string) (*authnClient.Tenant, error)
}

type authnApiClient struct {
	api *authnClient.APIClient
}

func NewAuthnClient(serviceURL string) (AuthnClient, error) {
	url, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	config := authnClient.NewConfiguration()
	config.Servers = []authnClient.ServerConfiguration{
		authnClient.ServerConfiguration{
			Url: url.String(),
		},
	}
	config.HTTPClient = tmw.UseClient(&http.Client{})

	return &authnApiClient{
		api: authnClient.NewAPIClient(config),
	}, nil
}

func (c *authnApiClient) GetTenant(claims tmw.TumblerClaims, tenantID string) (*authnClient.Tenant, error) {
	t, _, err := c.api.TenantsApi.GetTenant(claims.RequestContext(), tenantID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
