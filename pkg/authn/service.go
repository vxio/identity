package authn

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/session"
)

// AuthenticationApiServicer defines the api actions for the AuthenticationApi service
type AuthenticationService interface {
	LoginWithCredentials(req *http.Request, credentials client.Login, nonce string, ip string, photoURL *string) (*http.Cookie, *client.LoggedIn, error)
	RegisterWithCredentials(*http.Request, client.Register, string, string, bool) (*http.Cookie, *client.LoggedIn, error)
}

type authnService struct {
	log         logging.Logger
	credentials credentials.CredentialsService
	identities  identities.Service
	token       session.TokenService
	invites     invites.InvitesService
}

// NewAuthnService - Creates a default service that handles the registration and login
func NewAuthnService(
	log logging.Logger,
	credentials credentials.CredentialsService,
	identities identities.Service,
	token session.TokenService,
	invites invites.InvitesService,
) AuthenticationService {
	return &authnService{
		log:         log,
		credentials: credentials,
		identities:  identities,
		token:       token,
		invites:     invites,
	}
}

// RegisterWithCredentials - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (s *authnService) RegisterWithCredentials(req *http.Request, register client.Register, nonce string, ip string, isSignup bool) (*http.Cookie, *client.LoggedIn, error) {
	logCtx := s.log.WithMap(map[string]string{
		"tenant_id":     register.TenantID,
		"credential_id": register.CredentialID,
		"email":         register.Email,
		"ip":            ip,
		"is_sign_up":    fmt.Sprintf("%t", isSignup),
	})

	var invite *client.Invite = nil

	if !isSignup {
		invite, err := s.invites.Redeem(register.InviteCode)
		if err != nil {
			return nil, nil, logCtx.Error().LogErrorF("Unable to redeem token", err)
		}

		// Guard against possible inconsistencies in the tenantID
		if register.TenantID != invite.TenantID {
			return nil, nil, logCtx.Error().LogErrorF("register TenantID and Invite TenantID don't match")
		}
	}

	// See if this credential is already registered with the tenant and if so bomb out.
	found, err := s.credentials.Exists(register.CredentialID, register.TenantID)
	if err != nil {
		return nil, nil, err
	} else if found {
		return nil, nil, logCtx.Error().LogErrorF("credential already registered with tenant")
	}

	// Create the identity so we can login with it and give the user access.
	identity, err := s.identities.Register(register, invite)
	if err != nil {
		return nil, nil, logCtx.Error().LogErrorF("Unable to register identity", err)
	}

	// Register the credentials with the new Identity created.
	creds, err := s.credentials.Register(identity.IdentityID, register.CredentialID, register.TenantID)
	if err != nil {
		return nil, nil, logCtx.Error().LogErrorF("Unable to register credential", err)
	}

	// Using the new creds create the login object to log the user in.
	login := client.Login{
		CredentialID: creds.CredentialID,
		TenantID:     identity.TenantID,
	}

	return s.LoginWithCredentials(req, login, nonce, ip, register.ImageUrl)
}

// LoginWithCredentials - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service receives a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (s *authnService) LoginWithCredentials(req *http.Request, login client.Login, nonce string, ip string, photoURL *string) (*http.Cookie, *client.LoggedIn, error) {
	logCtx := s.log.WithMap(map[string]string{
		"tenant_id":     login.TenantID,
		"credential_id": login.CredentialID,
		"ip":            ip,
	})

	// check if they exist in the credentials service and if its enabled.
	credential, err := s.credentials.Login(login, nonce, ip)
	if err != nil {
		return nil, nil, logCtx.Error().LogError("Failed login", err)
	}

	logCtx = logCtx.With(api.NewCredentialLogContext(credential))

	identity, err := s.identities.GetIdentityByID(credential.IdentityID)
	if err != nil {
		return nil, nil, logCtx.Error().LogError("Could not find identity", err)
	}

	if identity.TenantID != credential.TenantID {
		return nil, nil, logCtx.LogErrorF("guard triggered - identity and credential tenantID's don't match")
	}

	if photoURL != nil && identity.ImageUrl != photoURL {
		logCtx.Info().Log("OIDC image updated")
		logCtx.Info().WithKeyValue("Old ImageUrl", *identity.ImageUrl)
		logCtx.Info().WithKeyValue("New ImageUrl", *photoURL)
		identity.ImageUrl = photoURL
		s.identities.UpdateInsecure(identity)
	}

	logCtx = logCtx.With(api.NewIdentityLogContext(identity))

	session := session.Session{
		IdentityID:   uuid.MustParse(identity.IdentityID),
		TenantID:     uuid.MustParse(identity.TenantID),
		CredentialID: uuid.MustParse(credential.CredentialID),
	}

	cookie, err := s.token.GenerateCookie(req, session)
	if err != nil {
		return nil, nil, logCtx.Error().LogError("Unable to generate cookie", err)
	}

	loggedIn := client.LoggedIn{
		Jwt:          cookie.Value,
		CredentialID: credential.CredentialID,
		TenantID:     credential.TenantID,
		IdentityID:   identity.IdentityID,
		FirstName:    identity.FirstName,
		LastName:     identity.LastName,
		NickName:     identity.NickName,
		ImageUrl:     identity.ImageUrl,
	}

	return cookie, &loggedIn, nil
}
