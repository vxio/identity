package database

import (
	"database/sql"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	migsqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
	"github.com/markbates/pkger"
)

func RunMigrations(db *sql.DB, config DatabaseConfig) error {
	fmt.Println("Running Migrations")

	pkger.Include("/migrations/")

	driver, err := GetDriver(db, config)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"pkger:///migrations/",
		config.DatabaseName,
		driver,
	)
	if err != nil {
		fmt.Printf("Error running migration - %s", err.Error())
		return err
	}

	err = m.Up()
	switch err {
	case nil:
	case migrate.ErrNoChange:
		fmt.Println("Database already at version")
	default:
		fmt.Printf("Error running migrations - %s\n", err.Error())
		return err
	}

	fmt.Println("Migrations complete")

	return nil
}

func GetDriver(db *sql.DB, config DatabaseConfig) (database.Driver, error) {
	if config.MySql != nil {
		return MySqlDriver(db)
	} else if config.SqlLite != nil {
		return Sqlite3Driver(db)
	}

	return nil, fmt.Errorf("Database config not defined")
}

func MySqlDriver(db *sql.DB) (database.Driver, error) {
	return migmysql.WithInstance(db, &migmysql.Config{})
}

func Sqlite3Driver(db *sql.DB) (database.Driver, error) {
	return migsqlite3.WithInstance(db, &migsqlite3.Config{})
}
