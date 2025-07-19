package store

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// ConnectDatabase attempts to connect to PostgreSQL first, then falls back to SQLite
func ConnectDatabase(dbURL, sqlitePath string) (*sqlx.DB, bool, error) {
	// Try PostgreSQL first
	if dbURL != "" {
		db, err := sqlx.Connect("postgres", dbURL)
		if err == nil {
			slog.Info("Connected to PostgreSQL database")
			return db, false, nil
		}
		slog.Warn("Failed to connect to PostgreSQL, trying SQLite fallback", "error", err)
	}

	// Fallback to SQLite
	return connectSQLite(sqlitePath)
}

// connectSQLite creates a SQLite database connection
func connectSQLite(sqlitePath string) (*sqlx.DB, bool, error) {
	// Ensure the directory exists
	dir := filepath.Dir(sqlitePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, false, fmt.Errorf("failed to create directory for SQLite database: %w", err)
	}

	// Connect to SQLite
	db, err := sqlx.Connect("sqlite", sqlitePath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, false, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	slog.Info("Connected to SQLite database", "path", sqlitePath)
	return db, true, nil
}
