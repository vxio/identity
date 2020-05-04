package main

import (
	"github.com/moov-io/identity/pkg/database"
)

type Config struct {
	Http     HttpConfig
	Admin    HttpConfig
	Database database.DatabaseConfig
}

type HttpConfig struct {
	Bind BindAddress
}

type BindAddress struct {
	Address string
}
