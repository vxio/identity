package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	//"github.com/moov-io/paygate/pkg/util"

	"github.com/go-kit/kit/log"
	"github.com/lopezator/migrator"
)

// Type returns a string for which database to be used.
func Type() string {
	dbType := "sqlite"
	if os.Getenv("DATABASE_TYPE") != "" {
		dbType = os.Getenv("DATABASE_TYPE")
	}
	return dbType
}

// New establishes a database connection according to the type and environmental
// variables for that specific database.
func New(ctx context.Context, logger log.Logger, _type string) (*sql.DB, error) {
	logger.Log("database", fmt.Sprintf("looking for %s database provider", _type))
	switch strings.ToLower(_type) {
	case "sqlite":
		return sqliteConnection(logger, getSqlitePath()).Connect(ctx)
	case "mysql":
		return mysqlConnection(logger, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_ADDRESS"), os.Getenv("MYSQL_DATABASE")).Connect(ctx)
	}
	return nil, fmt.Errorf("unknown database type %q", _type)
}

func execsql(name, raw string) *migrator.MigrationNoTx {
	return &migrator.MigrationNoTx{
		Name: name,
		Func: func(db *sql.DB) error {
			_, err := db.Exec(raw)
			return err
		},
	}
}

// UniqueViolation returns true when the provided error matches a database error
// for duplicate entries (violating a unique table constraint).
func UniqueViolation(err error) bool {
	return MySQLUniqueViolation(err) || SqliteUniqueViolation(err)
}
