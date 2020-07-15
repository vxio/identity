package api

import (
	"net/http"

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

// InternalApiRouter defines the required methods for binding the api requests to a responses for the InternalApi
// The InternalApiRouter implementation should parse necessary information from the http request,
// pass the data to a InternalApiServicer to perform the required actions, then write the service results to the http response.
type InternalApiRouter interface {
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
	DisableCredentials(tmw.TumblerClaims, string, string) (*Credential, error)
	ListCredentials(tmw.TumblerClaims, string) ([]Credential, error)
}

// IdentitiesApiServicer defines the api actions for the IdentitiesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type IdentitiesApiServicer interface {
	DisableIdentity(tmw.TumblerClaims, string) error
	GetIdentity(tmw.TumblerClaims, string) (*Identity, error)
	ListIdentities(tmw.TumblerClaims, string) ([]Identity, error)
	UpdateIdentity(tmw.TumblerClaims, string, UpdateIdentity) (*Identity, error)
}

// InternalApiServicer defines the api actions for the InternalApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type InternalApiServicer interface {
	LoginWithCredentials(*http.Request, Login, string, string) (*http.Cookie, error)
	RegisterWithCredentials(*http.Request, Register, string, string) (*http.Cookie, error)
}

// InvitesApiServicer defines the api actions for the InvitesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type InvitesApiServicer interface {
	DisableInvite(tmw.TumblerClaims, string) error
	ListInvites(tmw.TumblerClaims) ([]Invite, error)
	SendInvite(tmw.TumblerClaims, SendInvite) (*Invite, string, error)
	Redeem(code string) (*Invite, error)
}
