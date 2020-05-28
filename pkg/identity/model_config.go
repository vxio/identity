package identity

import (
	"github.com/moov-io/identity/pkg/authn"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/session"
)

// Config defines all the configuration for the app
type Config struct {
	Servers  ServerConfig
	Database database.DatabaseConfig
	//Keys           KeysConfig
	Gateway        gateway.Config
	Authentication authn.Config
	Session        session.Config
	Notifications  notifications.NotificationsConfig
	Invites        invites.Config
}

// ServerConfig - Groups all the http configs for the servers and ports that get opened.
type ServerConfig struct {
	Public HTTPConfig
	Admin  HTTPConfig
}

// HTTPConfig configuration for running an http server
type HTTPConfig struct {
	Bind BindAddress
}

// BindAddress specifies where the http server should bind to.
type BindAddress struct {
	Address string
}
