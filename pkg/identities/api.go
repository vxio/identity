package identities

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/zerotrust"
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

// DisableIdentity - Disable an identity. Its left around for historical reporting
func (c *Controller) DisableIdentity(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		err := c.service.DisableIdentity(session, identityID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(204)
	})
}

// GetIdentity - List identities and associates userId
func (c *Controller) GetIdentity(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		result, err := c.service.GetIdentity(session, identityID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}

// ListIdentities - List identities and associates userId
func (c *Controller) ListIdentities(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
		query := r.URL.Query()
		orgID := query.Get("orgID")
		result, err := c.service.ListIdentities(session, orgID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}

// UpdateIdentity - Update a specific Identity
func (c *Controller) UpdateIdentity(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
		params := mux.Vars(r)
		identityID := params["identityID"]
		identity := &api.UpdateIdentity{}
		if err := json.NewDecoder(r.Body).Decode(&identity); err != nil {
			w.WriteHeader(500)
			return
		}

		result, err := c.service.UpdateIdentity(session, identityID, *identity)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}
