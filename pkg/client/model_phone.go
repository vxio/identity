/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client

// Phone Phone number
type Phone struct {
	// UUID v4
	IdentityID string `json:"identityID,omitempty"`
	// UUID v4
	PhoneID   string `json:"phoneID,omitempty"`
	Number    string `json:"number,omitempty"`
	Validated bool   `json:"validated,omitempty"`
	Type      string `json:"type,omitempty"`
}
