package database

import (
	"database/sql"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
	"github.com/markbates/pkger"
)

func RunMigrations(db *sql.DB) error {
	fmt.Println("Running Migrations")

	pkger.Include("/migrations/")

	driver, err := migratedb.WithInstance(db, &migratedb.Config{})
	if err != nil {
		fmt.Printf("Error setting up migration - %s", err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"pkger:///migrations/",
		"identity",
		driver,
	)
	if err != nil {
		fmt.Printf("Error running migration - %s", err.Error())
		return nil
	}

	if err := m.Up(); err != nil {
		fmt.Printf("Error running migrations - %s\n", err.Error())
	}

	fmt.Println("Migrations complete")

	return nil
}
