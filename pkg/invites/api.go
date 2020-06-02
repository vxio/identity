package invites

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/gateway"
)

// A Controller binds http requests to an api service and writes the service results to the http response
type Controller struct {
	service api.InvitesApiServicer
}

// NewInvitesController creates a default api controller
func NewInvitesController(s api.InvitesApiServicer) api.Router {
	return &Controller{service: s}
}

// Routes returns all of the api route for the InvitesApiController
func (c *Controller) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "DeleteInvite",
			Method:      strings.ToUpper("Delete"),
			Pattern:     "/invites/{inviteID}",
			HandlerFunc: c.DeleteInvite,
		},
		{
			Name:        "ListInvites",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/invites",
			HandlerFunc: c.ListInvites,
		},
		{
			Name:        "SendInvite",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/invites",
			HandlerFunc: c.SendInvite,
		},
	}
}

// DeleteInvite - Delete an invite that was sent and invalidate the token.
func (c *Controller) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	gateway.WithSession(w, r, func(session gateway.Session) {
		params := mux.Vars(r)
		inviteID := params["inviteID"]
		err := c.service.DisableInvite(session, inviteID)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(204)
	})
}

// ListInvites - List outstanding invites
func (c *Controller) ListInvites(w http.ResponseWriter, r *http.Request) {
	gateway.WithSession(w, r, func(session gateway.Session) {
		//query := r.URL.Query()
		//orgID := query.Get("orgID")
		result, err := c.service.ListInvites(session)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}

// SendInvite - Send an email invite to a new user
func (c *Controller) SendInvite(w http.ResponseWriter, r *http.Request) {
	gateway.WithSession(w, r, func(session gateway.Session) {
		invite := &api.SendInvite{}
		if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
			w.WriteHeader(500)
			return
		}

		result, _, err := c.service.SendInvite(session, *invite)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		api.EncodeJSONResponse(result, nil, w)
	})
}
