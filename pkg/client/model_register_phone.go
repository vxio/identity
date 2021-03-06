/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client

// RegisterPhone Phone number
type RegisterPhone struct {
	Number string `json:"number,omitempty"`
	Type   string `json:"type,omitempty"`
}
