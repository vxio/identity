/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package authn

import (
	"errors"
	"net/http"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
)

// InternalService is a service that implents the logic for the InternalApiServicer
// This service should implement the business logic for every endpoint for the InternalApi API.
// Include any external packages or services that will be required by this service.
type InternalService struct {
	credentials credentials.CredentialsService
	identities  identities.IdentitiesService
	token       SessionService
	invites     invites.InvitesService
}

// NewInternalService creates a default api service
func NewAuthnService(
	credentials credentials.CredentialsService,
	identities identities.IdentitiesService,
	token SessionService,
) api.InternalApiServicer {
	return &InternalService{
		credentials: credentials,
		identities:  identities,
		token:       token,
	}
}

// RegisterWithCredentials - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (s *InternalService) RegisterWithCredentials(register api.Register) (*http.Cookie, error) {
	invite, err := s.invites.Redeem(register.InviteCode)
	if err != nil {
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

	return s.LoginWithCredentials(login)
}

// LoginWithCredentials - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service recieves a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (s *InternalService) LoginWithCredentials(login api.Login) (*http.Cookie, error) {
	// check if they exist in the credentials service and if its enabled.
	loggedIn, err := s.credentials.Login(login)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	// @TODO generate token with the ID's
	cookie, err := s.token.GenerateCookie(loggedIn.IdentityID)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}
