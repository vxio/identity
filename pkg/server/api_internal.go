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
	// None of these endpoints use this but the api generator adds it so if we need to use it again enable in the .openapi-generator-ignore
	// "github.com/gorilla/mux"
)

// A InternalApiController binds http requests to an api service and writes the service results to the http response
type InternalApiController struct {
	service InternalApiServicer
}

// NewInternalApiController creates a default api controller
func NewInternalApiController(s InternalApiServicer) Router {
	return &InternalApiController{service: s}
}

// Routes returns all of the api route for the InternalApiController
func (c *InternalApiController) Routes() Routes {
	return Routes{
		{
			"LoginPost",
			strings.ToUpper("Post"),
			"/login",
			c.LoginPost,
		},
		{
			"RegisterPost",
			strings.ToUpper("Post"),
			"/register",
			c.RegisterPost,
		},
	}
}

// LoginPost - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service recieves a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (c *InternalApiController) LoginPost(w http.ResponseWriter, r *http.Request) {
	login := &Login{}
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.LoginPost(*login)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// RegisterPost - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (c *InternalApiController) RegisterPost(w http.ResponseWriter, r *http.Request) {
	register := &Register{}
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := c.service.RegisterPost(*register)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
