package main

import (
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/jwks"
)

type Config struct {
	Http           HttpConfig
	Admin          HttpConfig
	Database       database.DatabaseConfig
	Authentication AuthenticationConfig
}

type HttpConfig struct {
	Bind BindAddress
}

type BindAddress struct {
	Address string
}

type AuthenticationConfig struct {
	Backchannel  jwks.JwksConfig
	Frontchannel jwks.JwksConfig
}
