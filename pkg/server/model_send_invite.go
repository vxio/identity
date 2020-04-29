/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform. 
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identityserver

// SendInvite - Describes an invite that was sent to a user to join.
type SendInvite struct {

	// Email Address
	Email string `json:"email,omitempty"`
}
