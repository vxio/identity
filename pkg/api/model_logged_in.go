/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform. 
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identityserver

// LoggedIn - User has logged in and is being given a token to proof identity
type LoggedIn struct {

	// JWT token that provides authentication of identity
	Jwt string `json:"jwt,omitempty"`
}