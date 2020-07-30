package authn

import (
	"net/http"
	"net/url"

	authnClient "github.com/moov-io/authn/pkg/client"
	"github.com/moov-io/identity/pkg/logging"

	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

type AuthnClient interface {
	GetTenant(claims tmw.TumblerClaims, tenantID string) (*authnClient.Tenant, error)
}

type authnApiClient struct {
	logger logging.Logger
	api    *authnClient.APIClient
}

func NewAuthnClient(logger logging.Logger, serviceURL string) (AuthnClient, error) {
	_, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	config := authnClient.NewConfiguration()
	config.BasePath = serviceURL
	config.Servers = []authnClient.ServerConfiguration{{Url: serviceURL}}
	config.HTTPClient = tmw.UseClient(&http.Client{})

	logger.WithKeyValue("base_path", serviceURL).Info().Log("Instantiated new Authn client")

	return &authnApiClient{
		logger: logger,
		api:    authnClient.NewAPIClient(config),
	}, nil
}

func (c *authnApiClient) GetTenant(claims tmw.TumblerClaims, tenantID string) (*authnClient.Tenant, error) {
	t, _, err := c.api.TenantsApi.GetTenant(claims.RequestContext(), tenantID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
