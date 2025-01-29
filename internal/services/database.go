package services

import (
	"database/sql"
	"log"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// SaveSwiftCodesToDatabase zapisuje dane SWIFT do bazy danych
func SaveSwiftCodesToDatabase(db *sql.DB, swiftCodes []models.SwiftCode) error {
	for _, code := range swiftCodes {
		query := `INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

		var id int64
		err := db.QueryRow(query, code.SwiftCode, code.BankName, code.Address, code.CountryISO2, code.CountryName, code.IsHeadquarter, code.HeadquarterID).Scan(&id)
		if err != nil {
			log.Printf("Błąd zapisu SWIFT code %s: %v", code.SwiftCode, err)
			return err
		}
	}

	log.Println("Wszystkie SWIFT codes zapisane w bazie.")
	return nil
}
