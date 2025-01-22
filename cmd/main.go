package main

import (
	"log"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/services"
)

func main() {
	// Dane do połączenia z bazą danych
	dsn := "host=localhost user=swiftuser password=mikus123 dbname=swift port=5432 sslmode=disable"

	// Inicjalizacja połączenia z bazą danych
	db.InitDatabase(dsn)
	defer db.CloseDatabase() // Zamknięcie połączenia po zakończeniu działania programu

	// Migracja bazy danych
	db.MigrateDatabase()

	// Parsowanie pliku CSV
	filePath := "data/swift_codes.csv"
	swiftCodes, err := services.ParseSwiftCodes(filePath)
	if err != nil {
		log.Fatalf("Error parsing SWIFT codes: %v", err)
	}

	// Zapisanie danych do bazy
	if err := services.SaveSwiftCodesToDatabase(swiftCodes); err != nil {
		log.Fatalf("Error saving SWIFT codes to database: %v", err)
	}

	log.Println("Dane zostały pomyślnie zapisane w bazie!")
}
