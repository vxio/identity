package credentials

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// A CredentialsApiController binds http requests to an api service and writes the service results to the http response
type credentialsApiController struct {
	service CredentialsService
}

// NewCredentialsApiController creates a default api controller
func NewCredentialsApiController(s CredentialsService) api.Router {
	return &credentialsApiController{service: s}
}

// Routes returns all of the api route for the CredentialsApiController
func (c *credentialsApiController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "DisableCredentials",
			Method:      strings.ToUpper("Delete"),
			Pattern:     "/identities/{identityID}/credentials/{credentialID}",
			HandlerFunc: c.DisableCredentials,
		},
		{
			Name:        "ListCredentials",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/identities/{identityID}/credentials",
			HandlerFunc: c.ListCredentials,
		},
	}
}

// DisableCredentials - Disables a credential so it can't be used anymore to login
func (c *credentialsApiController) DisableCredentials(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		credentialID := params["credentialID"]
		_, err := c.service.DisableCredentials(claims, identityID, credentialID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(404)
			default:
				w.WriteHeader(500)
			}

			return
		}

		w.WriteHeader(204)
	})
}

// ListCredentials - List the credentials this user has used.
func (c *credentialsApiController) ListCredentials(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		result, err := c.service.ListCredentials(claims, identityID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}
