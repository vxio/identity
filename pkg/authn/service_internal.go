package authn

import (
	"net/http"

	"github.com/google/uuid"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/session"
)

type authnService struct {
	log         logging.Logger
	credentials credentials.CredentialsService
	identities  identities.Service
	token       session.SessionService
	invites     api.InvitesApiServicer
}

// NewAuthnService - Creates a default service that handles the registration and login
func NewAuthnService(
	log logging.Logger,
	credentials credentials.CredentialsService,
	identities identities.Service,
	token session.SessionService,
	invites api.InvitesApiServicer,
) api.InternalApiServicer {
	return &authnService{
		log:         log,
		credentials: credentials,
		identities:  identities,
		token:       token,
		invites:     invites,
	}
}

// RegisterWithCredentials - Register user based on OIDC credentials.  This is called by the OIDC client services we create to register the user with what  available information they have and obtain from the user.
func (s *authnService) RegisterWithCredentials(req *http.Request, register api.Register, nonce string, ip string) (*http.Cookie, error) {
	logCtx := s.log.WithMap(map[string]string{
		"tenant_id":         register.TenantID,
		"credential_id":     register.CredentialID,
		"email":             register.Email,
		"invite_code_trunc": register.InviteCode[0:5],
		"ip":                ip,
	})

	invite, err := s.invites.Redeem(register.InviteCode)
	if err != nil {
		return nil, logCtx.Error().LogErrorF("Unable to redeem token", err)
	}

	// Guard against possible inconsistencies in the tenantID
	if register.TenantID != invite.TenantID {
		return nil, logCtx.Error().LogErrorF("register TenantID and Invite TenantID don't match")
	}

	// Create the identity so we can login with it and give the user access.
	identity, err := s.identities.Register(*invite, register)
	if err != nil {
		return nil, logCtx.Error().LogErrorF("Unable to register identity", err)
	}

	// Register the credentials with the new Identity created.
	creds, err := s.credentials.Register(identity.IdentityID, register.CredentialID, register.TenantID)
	if err != nil {
		return nil, logCtx.Error().LogErrorF("Unable to register credential", err)
	}

	// Using the new creds create the login object to log the user in.
	login := api.Login{
		CredentialID: creds.CredentialID,
		TenantID:     identity.TenantID,
	}

	return s.LoginWithCredentials(req, login, nonce, ip)
}

// LoginWithCredentials - Complete a login via a OIDC. Once the OIDC client service has authenticated their identity the client service will call  this endpoint to record and finish the login to get their token to use the API.  If the client service receives a 404 they must send them to registration if its allowed per the client or check for an invite for authenticated users email before sending to registration.
func (s *authnService) LoginWithCredentials(req *http.Request, login client.Login, nonce string, ip string) (*http.Cookie, error) {
	logCtx := s.log.WithMap(map[string]string{
		"tenant_id":     login.TenantID,
		"credential_id": login.CredentialID,
		"ip":            ip,
	})

	// check if they exist in the credentials service and if its enabled.
	loggedIn, err := s.credentials.Login(login, nonce, ip)
	if err != nil {
		return nil, logCtx.Error().LogError("Failed login", err)
	}

	logCtx = logCtx.With(api.NewCredentialLogContext(loggedIn))

	identity, err := s.identities.GetIdentityByID(loggedIn.IdentityID)
	if err != nil {
		return nil, logCtx.Error().LogError("Could not find identity", err)
	}

	if identity.TenantID != loggedIn.TenantID {
		return nil, logCtx.LogErrorF("guard triggered - identity and credential tenantID's don't match")
	}

	logCtx = logCtx.With(api.NewIdentityLogContext(identity))

	session := session.Session{
		IdentityID:   uuid.MustParse(identity.IdentityID),
		TenantID:     uuid.MustParse(identity.TenantID),
		CredentialID: uuid.MustParse(loggedIn.CredentialID),
	}

	cookie, err := s.token.GenerateCookie(req, session)
	if err != nil {
		return nil, logCtx.Error().LogError("Unable to generate cookie", err)
	}

	return cookie, nil
}
