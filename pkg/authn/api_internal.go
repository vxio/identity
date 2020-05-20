package authn

import (
	"net/http"
	"strings"

	"github.com/go-kit/kit/log"
	api "github.com/moov-io/identity/pkg/api"
)

// A AuthnAPIController binds http requests to an api service and writes the service results to the http response
type AuthnAPIController struct {
	logger  log.Logger
	service api.InternalApiServicer
}

// NewAuthnAPIController creates a default api controller
func NewAuthnAPIController(logger log.Logger, s api.InternalApiServicer) api.Router {
	return &AuthnAPIController{service: s}
}

// Routes returns all of the api route for the InternalApiController
func (c *AuthnAPIController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "Authenticated",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/authenticated",
			HandlerFunc: c.Authenticated,
		},
		{
			Name:        "Register",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/register",
			HandlerFunc: c.Register,
		},
		{
			Name:        "SubmitRegistration",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/register",
			HandlerFunc: c.Register,
		},
	}
}

// Authenticated - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service recieves a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (c *AuthnAPIController) Authenticated(w http.ResponseWriter, r *http.Request) {

	session, ok := r.Context().Value("LoginSession").(*LoginSession)
	if !ok || session == nil {
		c.logger.Log("level", "error", "msg", "Unable to find LoginSession in context")
		w.WriteHeader(404)
		return
	}

	login := api.Login{
		Provider:  session.Provider,
		SubjectID: session.SubjectID,
	}

	result, err := c.service.LoginWithCredentials(login, session.State, session.IP)
	if err != nil {
		c.logger.Log("level", "error", "msg", "Not able to exchange login token for session token", "error", err.Error())
		w.WriteHeader(404)
		return
	}

	http.SetCookie(w, result)
	http.Redirect(w, r, c.service.LandingURL(), http.StatusFound)
}

// Register - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (c *AuthnAPIController) Register(w http.ResponseWriter, r *http.Request) {

	// Show registration page but we don't really have one yet... so lets jut register with what we do have...
	c.SubmitRegistration(w, r)
}

func (c *AuthnAPIController) SubmitRegistration(w http.ResponseWriter, r *http.Request) {

	session, ok := r.Context().Value("LoginSession").(*LoginSession)
	if !ok || session == nil {
		c.logger.Log("level", "error", "msg", "Unable to find LoginSession in context")
		w.WriteHeader(404)
		return
	}

	register := api.Register{
		Provider:  session.Provider,
		SubjectID: session.SubjectID,

		InviteCode: session.InviteCode,

		FirstName:  session.FirstName,
		MiddleName: session.MiddleName,
		LastName:   session.LastName,

		NickName: session.NickName,

		Suffix:    session.Suffix,
		BirthDate: session.BirthDate,

		Email: session.Email,

		Phones:    session.Phones,
		Addresses: session.Addresses,
	}

	result, err := c.service.RegisterWithCredentials(register, session.State, session.IP)
	if err != nil {
		c.logger.Log("level", "error", "msg", "Unable to RegisterWithCredentials", "error", err)
		w.WriteHeader(400)
		return
	}

	http.SetCookie(w, result)
	http.Redirect(w, r, c.service.LandingURL(), http.StatusFound)
}
