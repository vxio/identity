package webkeys

type WebKeysConfig struct {
	File *FileConfig
	HTTP *HttpConfig
}

type FileConfig struct {
	Paths []string
}

type HttpConfig struct {
	URLs []string
}
