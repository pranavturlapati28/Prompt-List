package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var connStr string

	instanceConnection := os.Getenv("INSTANCE_CONNECTION_NAME")
	
	if instanceConnection != "" {
		socketDir := "/cloudsql"
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")
		
		connStr = fmt.Sprintf(
			"user=%s password=%s database=%s host=%s/%s sslmode=disable",
			dbUser, dbPass, dbName, socketDir, instanceConnection,
		)
	} else {
		connStr = os.Getenv("DATABASE_URL")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✓ Connected to database")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}