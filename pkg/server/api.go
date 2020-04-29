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
	"net/http"
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
	DisableCredentials(string, string) (interface{}, error)
	ListCredentials(string) (interface{}, error)
}

// IdentitiesApiServicer defines the api actions for the IdentitiesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type IdentitiesApiServicer interface {
	DisableIdentity(string) (interface{}, error)
	GetIdentity(string) (interface{}, error)
	ListIdentities(string) (interface{}, error)
	UpdateIdentity(string, UpdateIdentity) (interface{}, error)
}

// InternalApiServicer defines the api actions for the InternalApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type InternalApiServicer interface {
	LoginWithCredentials(Login) (interface{}, error)
	RegisterWithCredentials(Register) (interface{}, error)
}

// InvitesApiServicer defines the api actions for the InvitesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type InvitesApiServicer interface {
	DeleteInvite(string) (interface{}, error)
	ListInvites(string) (interface{}, error)
	SendInvite(SendInvite) (interface{}, error)
}
