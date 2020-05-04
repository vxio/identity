package main

type Config struct {
	Http  HttpConfig
	Admin HttpConfig
}

type HttpConfig struct {
	Bind BindAddress
}

type BindAddress struct {
	Address string
}
