package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Sterownik PostgreSQL dla database/sql
)

var DB *sql.DB // Globalne połączenie z bazą danych

// InitDatabase inicjalizuje połączenie z bazą PostgreSQL
func InitDatabase(dsn string) {
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Błąd połączenia z bazą danych: %v", err)
	}

	// Sprawdzenie połączenia
	if err := DB.Ping(); err != nil {
		log.Fatalf("Baza danych nie odpowiada: %v", err)
	}

	log.Println("Połączenie z bazą danych nawiązane")
}

// MigrateDatabase tworzy tabelę swift_codes, jeśli nie istnieje
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
		headquarter_id INT REFERENCES swift_codes(id) ON DELETE SET NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Błąd migracji bazy danych: %v", err)
	}

	log.Println("Migracja bazy danych zakończona sukcesem")
}

// CloseDatabase zamyka połączenie z bazą danych
func CloseDatabase() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Błąd zamykania połączenia z bazą danych: %v", err)
	}

	log.Println("Połączenie z bazą danych zamknięte")
}
