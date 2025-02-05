package repositories

import (
	"database/sql"
	"log"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// SwiftCodeRepository obsługuje operacje na tabeli swift_codes
type SwiftCodeRepository struct {
	db *sql.DB
}

// Nowe repozytorium SwiftCode
func NewSwiftCodeRepository(db *sql.DB) *SwiftCodeRepository {
	return &SwiftCodeRepository{db: db}
}

// Pobiera kod SWIFT na podstawie wartości
func (r *SwiftCodeRepository) GetBySwiftCode(code string) (*models.SwiftCode, error) {
	query := `
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE swift_code = $1;`

	swift := &models.SwiftCode{}
	var address sql.NullString // Obsługa potencjalnego NULL w bazie

	err := r.db.QueryRow(query, code).Scan(&swift.ID, &swift.SwiftCode, &swift.BankName, &address,
		&swift.CountryISO2, &swift.CountryName, &swift.IsHeadquarter, &swift.HeadquarterID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Brak wyniku
		}
		log.Println("Błąd pobierania SwiftCode:", err)
		return nil, err
	}

	// Jeśli adres jest NULL, ustawiamy wartość "UNKNOWN"
	if address.Valid {
		swift.Address = address.String
	} else {
		swift.Address = "UNKNOWN"
	}

	return swift, nil
}

// Pobiera listę SWIFT codes dla danego kraju (ISO-2)
func (r *SwiftCodeRepository) GetByCountryISO2(countryISO2 string) ([]models.SwiftCode, error) {
	query := `
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE country_iso2 = $1;`

	rows, err := r.db.Query(query, countryISO2)
	if err != nil {
		log.Println("Błąd pobierania SWIFT codes dla kraju:", err)
		return nil, err
	}
	defer rows.Close()

	var swiftCodes []models.SwiftCode
	for rows.Next() {
		var swift models.SwiftCode
		var address sql.NullString // Obsługa potencjalnego NULL w bazie

		err := rows.Scan(&swift.ID, &swift.SwiftCode, &swift.BankName, &address,
			&swift.CountryISO2, &swift.CountryName, &swift.IsHeadquarter, &swift.HeadquarterID)
		if err != nil {
			log.Println("Błąd skanowania rekordu SWIFT:", err)
			return nil, err
		}

		// Jeśli adres jest NULL, ustawiamy wartość "UNKNOWN"
		if address.Valid {
			swift.Address = address.String
		} else {
			swift.Address = "UNKNOWN"
		}

		swiftCodes = append(swiftCodes, swift)
	}
	return swiftCodes, nil
}

// Usuwa kod SWIFT z bazy danych
func (r *SwiftCodeRepository) DeleteSwiftCode(code string) error {
	query := "DELETE FROM swift_codes WHERE swift_code = $1;"
	_, err := r.db.Exec(query, code)
	if err != nil {
		log.Println("Błąd usuwania SWIFT code:", err)
		return err
	}
	return nil
}

// Odłącza wszystkie branche od headquartera
func (r *SwiftCodeRepository) DetachBranchesFromHeadquarter(headquarterID int64) error {
	query := "UPDATE swift_codes SET headquarter_id = NULL WHERE headquarter_id = $1;"
	_, err := r.db.Exec(query, headquarterID)
	if err != nil {
		log.Println("Błąd odłączania branchy:", err)
		return err
	}
	return nil
}

func (r *SwiftCodeRepository) InsertSwiftCode(swift *models.SwiftCode) error {
	query := `INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	return r.db.QueryRow(query, swift.SwiftCode, swift.BankName, swift.Address,
		swift.CountryISO2, swift.CountryName, swift.IsHeadquarter, swift.HeadquarterID).
		Scan(&swift.ID)
}

func (r *SwiftCodeRepository) GetBranchesByHeadquarter(headquarterCode string) ([]models.SwiftCode, error) {
	var headquarterID int

	// Pierwsze zapytanie: Pobieramy ID headquarters na podstawie swift_code
	err := r.db.QueryRow("SELECT id FROM swift_codes WHERE swift_code = $1;", headquarterCode).Scan(&headquarterID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Headquarter not found for SWIFT code:", headquarterCode)
			return nil, nil // Jeśli nie znaleziono headquarters, zwracamy pustą listę
		}
		log.Println("Error retrieving headquarter ID:", err)
		return nil, err
	}

	// Drugie zapytanie: Pobieramy listę branchy
	query := `
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE headquarter_id = $1;`

	rows, err := r.db.Query(query, headquarterID)
	if err != nil {
		log.Println("Error retrieving branches:", err)
		return nil, err
	}
	defer rows.Close()

	var branches []models.SwiftCode
	for rows.Next() {
		var branch models.SwiftCode
		var address sql.NullString // Obsługa potencjalnego NULL w bazie

		err := rows.Scan(&branch.ID, &branch.SwiftCode, &branch.BankName, &address,
			&branch.CountryISO2, &branch.CountryName, &branch.IsHeadquarter, &branch.HeadquarterID)
		if err != nil {
			log.Println("Error scanning branch record:", err)
			return nil, err
		}

		// Jeśli adres jest NULL, ustawiamy wartość "UNKNOWN"
		if address.Valid {
			branch.Address = address.String
		} else {
			branch.Address = "UNKNOWN"
		}

		branches = append(branches, branch)
	}

	return branches, nil
}
