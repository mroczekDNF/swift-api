package db

import (
	"log"

	"github.com/mroczekDNF/swift-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Globalne połączenie z bazą danych

// InitDatabase inicjalizuje połączenie z bazą danych PostgreSQL
func InitDatabase(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
}

// MigrateDatabase automatycznie tworzy tabele w bazie danych
func MigrateDatabase() {
	if err := DB.AutoMigrate(&models.SwiftCode{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")
}

// CloseDatabase zamyka połączenie z bazą danych
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve database connection: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}

	log.Println("Database connection closed")
}
