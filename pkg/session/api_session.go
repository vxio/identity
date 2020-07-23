package session

import (
	"net/http"
	"strings"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

type sessionController struct {
	log        logging.Logger
	service    SessionService
	identities identities.Service
	stime      stime.TimeService
}

// NewsessionController - Router for the Who Am I api routes.
func NewSessionController(log logging.Logger, service SessionService, identities identities.Service, stime stime.TimeService) api.Router {
	return &sessionController{
		log:        log,
		service:    service,
		identities: identities,
		stime:      stime,
	}
}

// Routes returns all of the api route for the AuthenticationApiController
func (c *sessionController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "session",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/session",
			HandlerFunc: c.session,
		},
	}
}

// session - Responds back with information about the authenticated session
func (c *sessionController) session(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		identity, err := c.identities.GetIdentity(claims, claims.IdentityID.String())
		if err != nil {
			c.log.Info().Log("Unable to lookup identity")
			w.WriteHeader(404)
			return
		}

		details := client.SessionDetails{
			CredentialID: claims.CredentialID.String(),
			TenantID:     claims.TenantID.String(),
			IdentityID:   claims.IdentityID.String(),
			FirstName:    identity.FirstName,
			LastName:     identity.LastName,
			NickName:     identity.NickName,
			ExpiresIn:    c.stime.Now().Unix() - claims.Expiry.Time().Unix(),
		}

		api.EncodeJSONResponse(&details, nil, w)
	})
}
