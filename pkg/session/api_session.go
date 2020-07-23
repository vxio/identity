package session

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

type SessionController interface {
	AppendRoutes(router *mux.Router) *mux.Router
}

func NewSessionController(logger logging.Logger, identities *identities.Service, stime stime.TimeService) SessionController {
	return &sessionController{
		logger:     logger,
		identities: identities,
		stime:      stime,
	}
}

type sessionController struct {
	logger     logging.Logger
	identities *identities.Service
	stime      stime.TimeService
}

func (c sessionController) AppendRoutes(router *mux.Router) *mux.Router {

	router.
		Name("Identity.getSession").
		Methods("GET").
		Path("/session").
		HandlerFunc(c.getSessionHandler)

	return router
}

func (c *sessionController) jsonResponse(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(value)
}

func (c *sessionController) errorResponse(w http.ResponseWriter, err error) {
	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(404)
	default:
		w.WriteHeader(500)
	}
}

func (c *sessionController) getSessionHandler(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		identity, err := c.identities.GetIdentity(claims, claims.IdentityID.String())
		if err != nil {
			c.logger.Info().Log("Unable to lookup identity")
			c.errorResponse(w, err)
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

		c.jsonResponse(w, &details)
	})
}
