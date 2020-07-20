/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform. 
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client
// RegisterAddressErrors Address of the Identity
type RegisterAddressErrors struct {
	// Descriptive reason for failing validation
	Type *string `json:"type,omitempty"`
	// Descriptive reason for failing validation
	Address1 *string `json:"address1,omitempty"`
	// Descriptive reason for failing validation
	Address2 *string `json:"address2,omitempty"`
	// Descriptive reason for failing validation
	City *string `json:"city,omitempty"`
	// Descriptive reason for failing validation
	State *string `json:"state,omitempty"`
	// Descriptive reason for failing validation
	PostalCode *string `json:"postalCode,omitempty"`
	// Descriptive reason for failing validation
	Country *string `json:"country,omitempty"`
}
