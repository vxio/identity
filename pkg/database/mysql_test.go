package database

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestMySQL__basic(t *testing.T) {
	db := CreateTestMySQLDB(t)
	defer db.Close()

	if err := db.DB.Ping(); err != nil {
		t.Fatal(err)
	}

	// create a phony MySQL
	m := mysqlConnection(log.NewNopLogger(), "user", "pass", "127.0.0.1:3006", "db")

	ctx, cancelFunc := context.WithCancel(context.Background())

	conn, err := m.Connect(ctx)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if conn != nil || err == nil {
		t.Fatalf("conn=%#v expected error", conn)
	}

	cancelFunc()
}

func TestMySQLUniqueViolation(t *testing.T) {
	err := errors.New(`problem upserting depository="282f6ffcd9ba5b029afbf2b739ee826e22d9df3b", userId="f25f48968da47ef1adb5b6531a1c2197295678ce": Error 1062: Duplicate entry '282f6ffcd9ba5b029afbf2b739ee826e22d9df3b' for key 'PRIMARY'`)
	if !UniqueViolation(err) {
		t.Error("should have matched unique violation")
	}
}
