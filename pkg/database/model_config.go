package database

type DatabaseConfig struct {
	MySql   *MySqlConfig
	SqlLite *SqlLiteConfig
}

type MySqlConfig struct {
	Address  string
	User     string
	Password string
	Database string
}

type SqlLiteConfig struct {
	Path string
}
