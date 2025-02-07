package integration

import (
	"log"
	"testing"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/services"
	"github.com/stretchr/testify/assert"
)

const testDBURL = "postgres://testuser:testpassword@localhost:5433/swift_test_db?sslmode=disable"

// SetupTestDatabase przygotowuje bazę danych do testów integracyjnych
func SetupTestDatabase(t *testing.T) {
	// Inicjalizacja połączenia do bazy danych
	db.InitDatabase(testDBURL)

	// Jeżeli test się nie powiedzie, zamykamy połączenie
	t.Cleanup(func() {
		if t.Failed() {
			log.Println("Zamykanie bazy danych po nieudanym teście")
		}
		db.CloseDatabase()
	})

	// Migracja struktury bazy danych
	db.MigrateDatabase()

	// Wczytanie danych testowych
	records, err := services.ParseSwiftCodes("../../data/test_data.csv")
	assert.NoError(t, err, "Błąd parsowania danych testowych")

	err = services.SaveSwiftCodesToDatabase(db.DB, records)
	assert.NoError(t, err, "Błąd zapisu danych testowych")

	log.Println("Baza testowa gotowa do testów")
}

// CleanupTestDatabase czyści dane z testowej bazy danych po zakończeniu testu
func CleanupTestDatabase(t *testing.T) {
	_, err := db.DB.Exec("TRUNCATE TABLE swift_codes RESTART IDENTITY CASCADE;")
	assert.NoError(t, err, "Błąd podczas czyszczenia bazy danych")
	log.Println("Baza testowa wyczyszczona")
}
