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

// Invite Describes an invite that was sent to a user to join.
type Invite struct {
	// UUID v4
	InviteID string `json:"inviteID,omitempty"`
	// UUID v4
	TenantID string `json:"tenantID,omitempty"`
	// Email Address
	Email string `json:"email,omitempty"`
	// UUID v4
	InvitedBy  string     `json:"invitedBy,omitempty"`
	InvitedOn  time.Time  `json:"invitedOn,omitempty"`
	RedeemedOn *time.Time `json:"redeemedOn,omitempty"`
	ExpiresOn  time.Time  `json:"expiresOn,omitempty"`
	DisabledOn *time.Time `json:"disabledOn,omitempty"`
	// UUID v4
	DisabledBy *string `json:"disabledBy,omitempty"`
}
