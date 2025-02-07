package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/routes"
	"github.com/mroczekDNF/swift-api/internal/services"
)

func main() {
	// Pobierz dane do połączenia z bazy danych z ENV
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Upewnij się, że wszystkie zmienne są ustawione
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatalf("Brakuje jednej lub więcej zmiennych środowiskowych: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	}

	// Tworzymy DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

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

	// Uruchomienie serwera
	r := routes.SetupRouter(db.DB)
	r.Run(":8080")
}
