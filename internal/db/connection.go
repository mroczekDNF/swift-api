package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for database/sql
)

var DB *sql.DB // Global database connection

// InitDatabase initializes the PostgreSQL database connection
func InitDatabase(dsn string) {
	var err error

	// Attempt connection up to 5 times
	for i := 0; i < 5; i++ {
		DB, err = sql.Open("pgx", dsn)
		if err == nil && DB.Ping() == nil {
			break
		}
		log.Printf("Database connection attempt failed (%d/5), retrying in 5 seconds...", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connection established")
}

// MigrateDatabase creates the swift_codes table if it does not exist
func MigrateDatabase() {
	query := `
	CREATE TABLE IF NOT EXISTS swift_codes (
		id SERIAL PRIMARY KEY,
		swift_code VARCHAR(11) UNIQUE NOT NULL,
		bank_name TEXT NOT NULL,
		address TEXT,
		country_iso2 CHAR(2) NOT NULL,
		country_name TEXT NOT NULL,
		is_headquarter BOOLEAN NOT NULL,
		headquarter_id INT
	);

	-- Add an index on headquarter_id for faster branch lookups
	CREATE INDEX IF NOT EXISTS idx_headquarter_id ON swift_codes (headquarter_id);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Database migration error: %v", err)
	}

	log.Println("Database migration completed successfully")
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing database connection: %v", err)
	}

	log.Println("Database connection closed")
}

// IsTableEmpty checks if a table is empty
func IsTableEmpty(tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
