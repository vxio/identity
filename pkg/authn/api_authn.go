package authn

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	log "github.com/moov-io/identity/pkg/logging"
)

// authnAPIController - Controller for the AuthN verification routes.
type authnAPIController struct {
	logger  log.Logger
	service AuthenticationService
}

// NewAuthnAPIController creates a default api controller
func NewAuthnAPIController(logger log.Logger, s AuthenticationService) api.Router {
	return &authnAPIController{logger: logger, service: s}
}

// Routes returns all of the api route for the AuthenticationApiController
func (c *authnAPIController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "Authenticated",
			Method:      strings.ToUpper("Post"),
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
	WithLoginSessionFromRequest(c.logger, w, r, []string{"authenticate", "finished"}, func(session LoginSession) {
		DeleteAuthnCookie(w)

		// Validation the session
		if err := validation.ValidateStruct(&session,
			validation.Field(&session.CredentialID, validation.Required, is.UUID),
			validation.Field(&session.TenantID, validation.Required, is.UUID),
			validation.Field(&session.State, validation.Required),
			validation.Field(&session.IP, validation.Required, is.IP),
		); err != nil {
			c.logger.Error().LogError("session validate failed", err)
			w.WriteHeader(404)
			return
		}

		login := client.Login{
			CredentialID: session.CredentialID,
			TenantID:     session.TenantID,
		}

		cookie, loggedIn, err := c.service.LoginWithCredentials(r, login, session.State, session.IP, session.ImageUrl)
		if err != nil {
			c.logger.Error().LogError("Not able to exchange login token for session token", err)
			w.WriteHeader(404)
			return
		}

		http.SetCookie(w, cookie)
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
	WithLoginSessionFromRequest(c.logger, w, r, []string{"register", "finished"}, func(session LoginSession) {

		// Validation the session
		if err := validation.ValidateStruct(&session,
			validation.Field(&session.CredentialID, validation.Required, is.UUID),
			validation.Field(&session.State, validation.Required),
			validation.Field(&session.IP, validation.Required, is.IP),
		); err != nil {
			c.logger.Error().LogError("unable to validate session", err)
			w.WriteHeader(404)
			return
		}

		// Going to overwrite or use what they've already sent.
		registration := &session.Register

		if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
			w.WriteHeader(400)
			return
		}

		// Validate the registration
		if err := registration.Validate(); err != nil {
			s := http.StatusBadRequest
			_ = api.EncodeJSONResponse(err, &s, w)
			return
		}

		// Check if this is a signup so we don't force the invite code lookup
		isSignup := false
		for _, v := range session.Scopes {
			if v == "signup" {
				isSignup = true
			}
		}

		DeleteAuthnCookie(w)

		cookie, loggedIn, err := c.service.RegisterWithCredentials(r, *registration, session.State, session.IP, isSignup)
		if err != nil {
			c.logger.Error().LogError("Unable to RegisterWithCredentials", err)
			w.WriteHeader(404)
			return
		}

		http.SetCookie(w, cookie)
		api.EncodeJSONResponse(loggedIn, nil, w)
	})
}
