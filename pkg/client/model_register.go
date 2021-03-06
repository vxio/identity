/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client

// Register Request to register a user in the system
type Register struct {
	// UUID v4
	CredentialID string `json:"credentialID,omitempty"`
	// UUID v4
	TenantID   string  `json:"tenantID,omitempty"`
	InviteCode string  `json:"inviteCode,omitempty"`
	FirstName  string  `json:"firstName,omitempty"`
	MiddleName string  `json:"middleName,omitempty"`
	LastName   string  `json:"lastName,omitempty"`
	NickName   *string `json:"nickName,omitempty"`
	ImageUrl   *string `json:"imageUrl,omitempty"`
	Suffix     *string `json:"suffix,omitempty"`
	BirthDate  *string `json:"birthDate,omitempty"`
	// Email Address
	Email     string            `json:"email,omitempty"`
	Phones    []RegisterPhone   `json:"phones,omitempty"`
	Addresses []RegisterAddress `json:"addresses,omitempty"`
}
