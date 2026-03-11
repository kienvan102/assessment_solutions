package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// NewInMemorySQLite initializes and returns a shared, in-memory SQLite database connection.
// The name parameter allows creating isolated memory databases that persist across the connection lifecycle.
func NewInMemorySQLite(name string) (*sql.DB, error) {
	// file:name?mode=memory&cache=shared creates a named in-memory database
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
	
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database '%s': %w", name, err)
	}

	// Ensure the connection is viable before returning
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database '%s': %w", name, err)
	}

	return db, nil
}
