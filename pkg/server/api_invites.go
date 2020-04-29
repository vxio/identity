/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identityserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A InvitesController binds http requests to an api service and writes the service results to the http response
type InvitesController struct {
	service InvitesApiServicer
}

// NewInvitesController creates a default api controller
func NewInvitesController(s InvitesApiServicer) Router {
	return &InvitesController{service: s}
}

// Routes returns all of the api route for the InvitesApiController
func (c *InvitesController) Routes() Routes {
	return Routes{
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
	params := mux.Vars(r)
	inviteID := params["inviteID"]
	result, err := c.service.DeleteInvite(inviteID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// ListInvites - List outstanding invites
func (c *InvitesController) ListInvites(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	orgID := query.Get("orgID")
	result, err := c.service.ListInvites(orgID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// SendInvite - Send an email invite to a new user
func (c *InvitesController) SendInvite(w http.ResponseWriter, r *http.Request) {
	invite := &Invite{}
	if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.SendInvite(*invite)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
