package identities

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// A Controller binds http requests to an api service and writes the service results to the http response
type Controller struct {
	service api.IdentitiesApiServicer
}

// NewIdentitiesController creates a default api controller
func NewIdentitiesController(s api.IdentitiesApiServicer) api.Router {
	return &Controller{service: s}
}

// Routes returns all of the api route for the IdentitiesApiController
func (c *Controller) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "DisableIdentity",
			Method:      strings.ToUpper("Delete"),
			Pattern:     "/identities/{identityID}",
			HandlerFunc: c.DisableIdentity,
		},
		{
			Name:        "GetIdentity",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/identities/{identityID}",
			HandlerFunc: c.GetIdentity,
		},
		{
			Name:        "ListIdentities",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/identities",
			HandlerFunc: c.ListIdentities,
		},
		{
			Name:        "UpdateIdentity",
			Method:      strings.ToUpper("Put"),
			Pattern:     "/identities/{identityID}",
			HandlerFunc: c.UpdateIdentity,
		},
	}
}

func errorHandling(w http.ResponseWriter, err error) {
	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(404)
	default:
		w.WriteHeader(500)
		return
	}
}

// DisableIdentity - Disable an identity. Its left around for historical reporting
func (c *Controller) DisableIdentity(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		err := c.service.DisableIdentity(claims, identityID)
		if err != nil {
			errorHandling(w, err)
			return
		}

		w.WriteHeader(204)
	})
}

// GetIdentity - List identities and associates userId
func (c *Controller) GetIdentity(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		result, err := c.service.GetIdentity(claims, identityID)
		if err != nil {
			errorHandling(w, err)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}

// ListIdentities - List identities and associates userId
func (c *Controller) ListIdentities(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		query := r.URL.Query()
		orgID := query.Get("orgID")
		result, err := c.service.ListIdentities(claims, orgID)
		if err != nil {
			errorHandling(w, err)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}

// UpdateIdentity - Update a specific Identity
func (c *Controller) UpdateIdentity(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		identity := &api.UpdateIdentity{}
		if err := json.NewDecoder(r.Body).Decode(&identity); err != nil {
			w.WriteHeader(400)
			return
		}

		result, err := c.service.UpdateIdentity(claims, identityID, *identity)
		if err != nil {
			errorHandling(w, err)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}
