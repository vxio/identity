package main

import (
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/jwks"
	"github.com/moov-io/identity/pkg/notifications"
)

//Config defines all the configuration for the app
type Config struct {
	HTTP           HTTPConfig
	Admin          HTTPConfig
	Database       database.DatabaseConfig
	Authentication AuthenticationConfig
	Notifications  notifications.NotificationsConfig
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
	Backchannel  jwks.JwksConfig
	Frontchannel jwks.JwksConfig
}
