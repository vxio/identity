package identity

import (
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/webkeys"
)

//IdentityConfig defines all the configuration for the app
type IdentityConfig struct {
	Servers        ServerConfig
	Database       database.DatabaseConfig
	Authentication AuthenticationConfig
	Notifications  notifications.NotificationsConfig
	Invites        invites.InvitesConfig
}

type ServerConfig struct {
	Public  HTTPConfig
	Private HTTPConfig
	Admin   HTTPConfig
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
type AuthenticationConfig struct {
	Backchannel  webkeys.WebKeysConfig
	Frontchannel webkeys.WebKeysConfig
}
