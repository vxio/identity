package webkeys

type WebKeysConfig struct {
	File *FileConfig
	HTTP *HttpConfig
}

type FileConfig struct {
	Path string
}

type HttpConfig struct {
	URL string
}
