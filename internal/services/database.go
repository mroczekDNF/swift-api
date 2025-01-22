package services

import (
	"log"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// SaveSwiftCodesToDatabase zapisuje dane SWIFT do bazy danych
func SaveSwiftCodesToDatabase(swiftCodes []models.SwiftCode) error {
	for _, code := range swiftCodes {
		// Próba zapisania rekordu do bazy
		if err := db.DB.Create(&code).Error; err != nil {
			log.Printf("Failed to save SWIFT code %s: %v", code.SwiftCode, err)
			return err
		}
	}
	log.Println("All SWIFT codes have been successfully saved to the database.")
	return nil
}
