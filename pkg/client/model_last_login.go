/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform. 
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client
import (
	"time"
)
// LastLogin Defines when and what credential was used for the last login 
type LastLogin struct {
	// UUID v4
	CredentialId string `json:"credentialId,omitempty"`
	On time.Time `json:"on,omitempty"`
}
