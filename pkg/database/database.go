package database

import (
	"context"
	"database/sql"
	"fmt"

	//"github.com/moov-io/paygate/pkg/util"

	"github.com/go-kit/kit/log"
	"github.com/lopezator/migrator"
)

// New establishes a database connection according to the type and environmental
// variables for that specific database.
func New(ctx context.Context, logger log.Logger, config DatabaseConfig) (*sql.DB, error) {
	if config.MySql != nil {
		return mysqlConnection(logger, config.MySql.User, config.MySql.Password, config.MySql.Address, config.MySql.Database).Connect(ctx)
	} else if config.SqlLite != nil {
		return sqliteConnection(logger, config.SqlLite.Path).Connect(ctx)
	}

	return nil, fmt.Errorf("Database config not defined")
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
