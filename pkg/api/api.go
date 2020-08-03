package api

import (
	"net/http"

	"github.com/moov-io/identity/pkg/client"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

// CredentialsApiRouter defines the required methods for binding the api requests to a responses for the CredentialsApi
// The CredentialsApiRouter implementation should parse necessary information from the http request,
// pass the data to a CredentialsApiServicer to perform the required actions, then write the service results to the http response.
type CredentialsApiRouter interface {
	DisableCredentials(http.ResponseWriter, *http.Request)
	ListCredentials(http.ResponseWriter, *http.Request)
}

// IdentitiesApiRouter defines the required methods for binding the api requests to a responses for the IdentitiesApi
// The IdentitiesApiRouter implementation should parse necessary information from the http request,
// pass the data to a IdentitiesApiServicer to perform the required actions, then write the service results to the http response.
type IdentitiesApiRouter interface {
	DisableIdentity(http.ResponseWriter, *http.Request)
	GetIdentity(http.ResponseWriter, *http.Request)
	ListIdentities(http.ResponseWriter, *http.Request)
	UpdateIdentity(http.ResponseWriter, *http.Request)
}

// AuthenticationApiRouter defines the required methods for binding the api requests to a responses for the AuthenticationApi
// The AuthenticationApiRouter implementation should parse necessary information from the http request,
// pass the data to a AuthenticationApiServicer to perform the required actions, then write the service results to the http response.
type AuthenticationApiRouter interface {
	LoginWithCredentials(http.ResponseWriter, *http.Request)
	RegisterWithCredentials(http.ResponseWriter, *http.Request)
}

// InvitesApiRouter defines the required methods for binding the api requests to a responses for the InvitesApi
// The InvitesApiRouter implementation should parse necessary information from the http request,
// pass the data to a InvitesApiServicer to perform the required actions, then write the service results to the http response.
type InvitesApiRouter interface {
	DeleteInvite(http.ResponseWriter, *http.Request)
	ListInvites(http.ResponseWriter, *http.Request)
	SendInvite(http.ResponseWriter, *http.Request)
}

// CredentialsApiServicer defines the api actions for the CredentialsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CredentialsApiServicer interface {
	DisableCredentials(tmw.TumblerClaims, string, string) (*client.Credential, error)
	ListCredentials(tmw.TumblerClaims, string) ([]client.Credential, error)
}

// IdentitiesApiServicer defines the api actions for the IdentitiesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type IdentitiesApiServicer interface {
	DisableIdentity(tmw.TumblerClaims, string) error
	GetIdentity(tmw.TumblerClaims, string) (*client.Identity, error)
	ListIdentities(tmw.TumblerClaims, string) ([]client.Identity, error)
	UpdateIdentity(tmw.TumblerClaims, string, client.UpdateIdentity) (*client.Identity, error)
}

// AuthenticationApiServicer defines the api actions for the AuthenticationApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AuthenticationApiServicer interface {
	LoginWithCredentials(*http.Request, client.Login, string, string) (*http.Cookie, *client.LoggedIn, error)
	RegisterWithCredentials(*http.Request, client.Register, string, string, bool) (*http.Cookie, *client.LoggedIn, error)
}

// InvitesApiServicer defines the api actions for the InvitesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type InvitesApiServicer interface {
	DisableInvite(tmw.TumblerClaims, string) error
	ListInvites(tmw.TumblerClaims) ([]client.Invite, error)
	SendInvite(tmw.TumblerClaims, client.SendInvite) (*client.Invite, string, error)
	Redeem(code string) (*client.Invite, error)
}
