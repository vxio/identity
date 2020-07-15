package authn

import (
	"encoding/json"
	"net/http"
	"strings"

	api "github.com/moov-io/identity/pkg/api"
	log "github.com/moov-io/identity/pkg/logging"
)

// authnAPIController - Controller for the AuthN verification routes.
type authnAPIController struct {
	logger  log.Logger
	service api.InternalApiServicer
}

// NewAuthnAPIController creates a default api controller
func NewAuthnAPIController(logger log.Logger, s api.InternalApiServicer) api.Router {
	return &authnAPIController{logger: logger, service: s}
}

// Routes returns all of the api route for the InternalApiController
func (c *authnAPIController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "Authenticated",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/authentication/authenticated",
			HandlerFunc: c.Authenticated,
		},
		{
			Name:        "Register",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/authentication/register",
			HandlerFunc: c.Register,
		},
		{
			Name:        "SubmitRegistration",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/authentication/register",
			HandlerFunc: c.SubmitRegistration,
		},
	}
}

// Authenticated - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service receives a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (c *authnAPIController) Authenticated(w http.ResponseWriter, r *http.Request) {
	WithLoginSessionFromRequest(c.logger, w, r, []string{"authenticate"}, func(session LoginSession) {
		DeleteAuthnCookie(w)

		login := api.Login{
			CredentialID: session.CredentialID,
			TenantID:     session.TenantID,
		}

		result, err := c.service.LoginWithCredentials(r, login, session.State, session.IP)
		if err != nil {
			c.logger.Error().LogError("Not able to exchange login token for session token", err)
			w.WriteHeader(404)
			return
		}

		loggedIn := api.LoggedIn{
			Jwt: result.Value,
		}

		http.SetCookie(w, result)
		api.EncodeJSONResponse(loggedIn, nil, w)
	})
}

// Register - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (c *authnAPIController) Register(w http.ResponseWriter, r *http.Request) {
	// Show registration page but we don't really have one yet... so lets jut register with what we do have...
	WithLoginSessionFromRequest(c.logger, w, r, []string{"register"}, func(session LoginSession) {
		api.EncodeJSONResponse(session.Register, nil, w)
	})
}

// SubmitRegistration - Finalizes the registration and handles all the user creation and first login
func (c *authnAPIController) SubmitRegistration(w http.ResponseWriter, r *http.Request) {
	WithLoginSessionFromRequest(c.logger, w, r, []string{"register"}, func(session LoginSession) {
		registration := &api.Register{}
		if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
			w.WriteHeader(400)
			return
		}

		DeleteAuthnCookie(w)

		result, err := c.service.RegisterWithCredentials(r, *registration, session.State, session.IP)
		if err != nil {
			c.logger.Error().LogError("Unable to RegisterWithCredentials", err)
			w.WriteHeader(400)
			return
		}

		loggedIn := api.LoggedIn{
			Jwt: result.Value,
		}

		http.SetCookie(w, result)
		api.EncodeJSONResponse(loggedIn, nil, w)
	})
}
