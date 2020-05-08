package webkeys

import "time"

type WebKeysConfig struct {
	File *FileConfig
	HTTP *HttpConfig

	Expiration time.Duration
}

type FileConfig struct {
	Path string
}

type HttpConfig struct {
	URL string
}
