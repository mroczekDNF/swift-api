package integration

import (
	"log"
	"testing"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/services"
	"github.com/stretchr/testify/assert"
)

const testDBURL = "postgres://testuser:testpassword@localhost:5433/swift_test_db?sslmode=disable"

func SetupTestDatabase(t *testing.T) {
	db.InitDatabase(testDBURL)

	t.Cleanup(func() {
		if t.Failed() {
			log.Println("Closing database connection after a failed test")
		}
		db.CloseDatabase()
	})

	db.MigrateDatabase()

	records, err := services.ParseSwiftCodes("../../data/test_data.csv")
	assert.NoError(t, err, "Error parsing test data")

	err = services.SaveSwiftCodesToDatabase(db.DB, records)
	assert.NoError(t, err, "Error saving test data to the database")
}

func CleanupTestDatabase(t *testing.T) {
	_, err := db.DB.Exec("TRUNCATE TABLE swift_codes RESTART IDENTITY CASCADE;")
	assert.NoError(t, err, "Error cleaning up the test database")
	log.Println("Test database cleaned up")
}
