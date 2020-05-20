package authn

import (
	"errors"
	"fmt"
	"net/http"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/identities"
)

// InternalService is a service that implents the logic for the InternalApiServicer
// This service should implement the business logic for every endpoint for the InternalApi API.
// Include any external packages or services that will be required by this service.
type internalService struct {
	credentials credentials.CredentialsService
	identities  identities.Service
	token       SessionService
	invites     api.InvitesApiServicer
	landingURL  string
}

// NewInternalService creates a default api service
func NewAuthnService(
	credentials credentials.CredentialsService,
	identities identities.Service,
	token SessionService,
	invites api.InvitesApiServicer,
	landingURL string,
) api.InternalApiServicer {
	return &internalService{
		credentials: credentials,
		identities:  identities,
		token:       token,
		invites:     invites,
		landingURL:  landingURL,
	}
}

// RegisterWithCredentials - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (s *internalService) RegisterWithCredentials(register api.Register, nonce string, ip string) (*http.Cookie, error) {
	invite, err := s.invites.Redeem(register.InviteCode)
	if err != nil {
		fmt.Println("Unable to redeem token", register.InviteCode, err)
		return nil, err
	}

	// Create the identity so we can login with it and give the user access.
	identity, err := s.identities.Register(*invite, register)
	if err != nil {
		return nil, err
	}

	// Register the credentials with the new Identity created.
	creds, err := s.credentials.Register(identity.IdentityID, register.Provider, register.SubjectID)
	if err != nil {
		return nil, err
	}

	// Using the new creds create the login object to log the user in.
	login := api.Login{
		Provider:  creds.Provider,
		SubjectID: creds.SubjectID,
	}

	return s.LoginWithCredentials(login, nonce, ip)
}

// LoginWithCredentials - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service recieves a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (s *internalService) LoginWithCredentials(login api.Login, nonce string, ip string) (*http.Cookie, error) {
	// check if they exist in the credentials service and if its enabled.
	loggedIn, err := s.credentials.Login(login, nonce, ip)
	if err != nil {
		fmt.Println("Failed login ", err.Error())
		return nil, errors.New("Unauthorized")
	}

	// @TODO generate token with the ID's
	cookie, err := s.token.GenerateCookie(loggedIn.IdentityID)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func (s *internalService) LandingURL() string {
	return s.landingURL
}
