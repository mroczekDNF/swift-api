package main

import (
	"log"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/routes"
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
		log.Fatalf("Błąd parsowania SWIFT codes: %v", err)
	}

	// Zapisanie danych do bazy
	if err := services.SaveSwiftCodesToDatabase(db.DB, swiftCodes); err != nil {
		log.Fatalf("Błąd zapisu SWIFT codes do bazy: %v", err)
	}

	log.Println("Dane zostały pomyślnie zapisane w bazie!")

	r := routes.SetupRouter(db.DB)
	r.Run(":8080")
}
