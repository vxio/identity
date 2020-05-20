package identity

import (
	"github.com/moov-io/identity/pkg/authn"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/webkeys"
)

//IdentityConfig defines all the configuration for the app
type IdentityConfig struct {
	Servers       ServerConfig
	Database      database.DatabaseConfig
	Keys          KeysConfig
	Session       authn.SessionConfig
	Notifications notifications.NotificationsConfig
	Invites       invites.Config
}

type ServerConfig struct {
	Public HTTPConfig
	Admin  HTTPConfig
}

//HTTPConfig configuration for running an http server
type HTTPConfig struct {
	Bind BindAddress
}

//BindAddress specifies where the http server should bind to.
type BindAddress struct {
	Address string
}

//AuthenticationConfig on where to get keys from.
//  Backchannel is for verifying what comes from the Gateway
//  Frontchannel is for creating the tokens sent to the customer.
type KeysConfig struct {
	AuthnPublic    webkeys.WebKeysConfig
	GatewayPublic  webkeys.WebKeysConfig
	SessionPublic  webkeys.WebKeysConfig
	SessionPrivate webkeys.WebKeysConfig
}
