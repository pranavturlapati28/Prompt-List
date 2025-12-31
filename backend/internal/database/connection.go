package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection pool
var DB *sql.DB

// Connect establishes a connection to PostgreSQL
func Connect(databaseURL string) error {
	var err error

	// Open a connection pool to the database
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection actually works
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure the connection pool
	DB.SetMaxOpenConns(25)                  // Maximum number of open connections
	DB.SetMaxIdleConns(5)                   // Maximum number of idle connections
	DB.SetConnMaxLifetime(5 * 60 * 1000000000) // 5 minutes in nanoseconds

	fmt.Println("✓ Connected to PostgreSQL")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
		fmt.Println("✓ Database connection closed")
	}
}