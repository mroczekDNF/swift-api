package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Sterownik PostgreSQL dla database/sql
)

var DB *sql.DB // Globalne połączenie z bazą danych

// InitDatabase inicjalizuje połączenie z bazą PostgreSQL
func InitDatabase(dsn string) {
	var err error

	// 5 prób połączenia
	for i := 0; i < 5; i++ {
		DB, err = sql.Open("pgx", dsn)
		if err == nil && DB.Ping() == nil {
			break
		}
		log.Printf("Próba połączenia z bazą danych nie powiodła się (%d/5), ponawianie za 5 sekund...", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Nie udało się połączyć z bazą danych: %v", err)
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
		headquarter_id INT
	);

	-- Dodanie indeksu na headquarter_id dla szybkiego wyszukiwania branchy
	CREATE INDEX IF NOT EXISTS idx_headquarter_id ON swift_codes (headquarter_id);
	`
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
