package session

import (
	"fmt"
	"net/http"
	"strings"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/identities"
)

type whoAmIController struct {
	service    SessionService
	identities identities.Service
}

// NewWhoAmIController - Router for the Who Am I api routes.
func NewWhoAmIController(service SessionService, identities identities.Service) api.Router {
	return &whoAmIController{
		service:    service,
		identities: identities,
	}
}

// Routes returns all of the api route for the InternalApiController
func (c *whoAmIController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "WhoAmI",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/whoami",
			HandlerFunc: c.WhoAmI,
		},
	}
}

// WhoAmI - Responds back with information about the authenticated session
func (c *whoAmIController) WhoAmI(w http.ResponseWriter, r *http.Request) {
	cookieSession, err := c.service.FromRequest(r)
	if err != nil {
		fmt.Println("Cookie session not found")
		w.WriteHeader(404)
		return
	}

	gatewaySession, err := gateway.SessionFromRequest(r)
	if err != nil {
		fmt.Println("Gateway session not found")
		w.WriteHeader(404)
		return
	}

	identity, err := c.identities.GetIdentity(*gatewaySession, gatewaySession.CallerID.String())
	if err != nil {
		fmt.Println("Unable to lookup identity")
		w.WriteHeader(404)
		return
	}

	type Output struct {
		Cookie   Session
		Gateway  gateway.Session
		Identity api.Identity
		XUser    string
		XTenant  string
	}

	output := Output{
		Cookie:   *cookieSession,
		Gateway:  *gatewaySession,
		Identity: *identity,
		XUser:    r.Header.Get("X-User"),
		XTenant:  r.Header.Get("X-Tenant"),
	}

	api.EncodeJSONResponse(output, nil, w)
}
