/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client

// RegisterAddress Address of the Identity
type RegisterAddress struct {
	Type       string  `json:"type,omitempty"`
	Address1   string  `json:"address1,omitempty"`
	Address2   *string `json:"address2,omitempty"`
	City       string  `json:"city,omitempty"`
	State      string  `json:"state,omitempty"`
	PostalCode string  `json:"postalCode,omitempty"`
	Country    string  `json:"country,omitempty"`
}
