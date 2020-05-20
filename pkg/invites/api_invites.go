package invites

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/zerotrust"
)

// A InvitesController binds http requests to an api service and writes the service results to the http response
type InvitesController struct {
	service api.InvitesApiServicer
}

// NewInvitesController creates a default api controller
func NewInvitesController(s api.InvitesApiServicer) api.Router {
	return &InvitesController{service: s}
}

// Routes returns all of the api route for the InvitesApiController
func (c *InvitesController) Routes() api.Routes {
	return api.Routes{
		{
			"DeleteInvite",
			strings.ToUpper("Delete"),
			"/invites/{inviteID}",
			c.DeleteInvite,
		},
		{
			"ListInvites",
			strings.ToUpper("Get"),
			"/invites",
			c.ListInvites,
		},
		{
			"SendInvite",
			strings.ToUpper("Post"),
			"/invites",
			c.SendInvite,
		},
	}
}

// DeleteInvite - Delete an invite that was sent and invalidate the token.
func (c *InvitesController) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
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
func (c *InvitesController) ListInvites(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
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
func (c *InvitesController) SendInvite(w http.ResponseWriter, r *http.Request) {
	zerotrust.WithSession(w, r, func(session zerotrust.Session) {
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
