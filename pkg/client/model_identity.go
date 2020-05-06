/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform. 
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identityclient
import (
	"time"
)
// Identity Properties of an Identity. These users will under-go KYC checks thus all the information 
type Identity struct {
	// UUID v4
	IdentityID string `json:"identityID,omitempty"`
	// UUID v4
	TenantID string `json:"tenantID,omitempty"`
	FirstName string `json:"firstName"`
	MiddleName string `json:"middleName,omitempty"`
	LastName string `json:"lastName"`
	NickName *string `json:"nickName,omitempty"`
	Suffix *string `json:"suffix,omitempty"`
	BirthDate string `json:"birthDate,omitempty"`
	Status string `json:"status,omitempty"`
	// Email Address
	Email string `json:"email"`
	// The user has verified they have access to this email
	EmailVerified bool `json:"emailVerified,omitempty"`
	Phones []Phone `json:"phones,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
	RegisteredOn time.Time `json:"registeredOn,omitempty"`
	LastLogin LastLogin `json:"lastLogin,omitempty"`
	DisabledOn *time.Time `json:"disabledOn,omitempty"`
	// UUID v4
	DisabledBy *string `json:"disabledBy,omitempty"`
	LastUpdatedOn time.Time `json:"lastUpdatedOn,omitempty"`
}
