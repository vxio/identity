package session

import (
	"net/http"
	"strings"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
)

type whoAmIController struct {
	log        logging.Logger
	service    SessionService
	identities identities.Service
}

// NewWhoAmIController - Router for the Who Am I api routes.
func NewWhoAmIController(log logging.Logger, service SessionService, identities identities.Service) api.Router {
	return &whoAmIController{
		log:        log,
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
		c.log.Info().Log("Cookie session not found")
		w.WriteHeader(404)
		return
	}

	gatewaySession, err := gateway.SessionFromRequest(r)
	if err != nil {
		c.log.Info().Log("Gateway session not found")
		w.WriteHeader(404)
		return
	}

	logCtx := c.log.With(gatewaySession)

	identity, err := c.identities.GetIdentity(*gatewaySession, gatewaySession.CallerID.String())
	if err != nil {
		logCtx.Info().Log("Unable to lookup identity")
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
