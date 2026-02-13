package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	db.SetMaxOpenConns(1) // Required for :memory: database

	if err := createTablesForDB(db); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return db
}

func TeardownTestDB(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Errorf("Failed to close test database: %v", err)
	}
}
