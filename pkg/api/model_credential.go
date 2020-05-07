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
	"time"
)

// Credential - Description of a successful OpenID connect credential 
type Credential struct {

	// UUID v4
	CredentialID string `json:"credentialID,omitempty"`

	// OIDC provider that was used to handle authentication of this user.
	Provider string `json:"provider,omitempty"`

	// ID of the remote OIDC server gives to this identity
	SubjectID string `json:"subjectID,omitempty"`

	// UUID v4
	IdentityID string `json:"identityID,omitempty"`

	CreatedOn time.Time `json:"createdOn,omitempty"`

	LastUsedOn time.Time `json:"lastUsedOn,omitempty"`

	DisabledOn *time.Time `json:"disabledOn,omitempty"`

	// UUID v4
	DisabledBy *string `json:"disabledBy,omitempty"`
}