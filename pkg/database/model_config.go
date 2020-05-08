package database

type DatabaseConfig struct {
	MySql        *MySqlConfig
	SqlLite      *SqlLiteConfig
	DatabaseName string
}

type MySqlConfig struct {
	Address  string
	User     string
	Password string
}

type SqlLiteConfig struct {
	Path string
}
