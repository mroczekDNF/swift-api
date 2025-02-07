package repositories

import (
	"database/sql"
	"log"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// SwiftCodeRepositoryInterface opisuje metody repozytorium
type SwiftCodeRepositoryInterface interface {
	GetBySwiftCode(code string) (*models.SwiftCode, error)
	GetByCountryISO2(countryISO2 string) ([]models.SwiftCode, error)
	DeleteSwiftCode(code string) error
	DetachBranchesFromHeadquarter(headquarterID int64) error
	InsertSwiftCode(swift *models.SwiftCode) error
	GetBranchesByHeadquarter(headquarterCode string) ([]models.SwiftCode, error)
}

// SwiftCodeRepository obsługuje operacje na tabeli swift_codes
type SwiftCodeRepository struct {
	db *sql.DB
}

// NewSwiftCodeRepository tworzy nowe repozytorium SwiftCode
func NewSwiftCodeRepository(db *sql.DB) *SwiftCodeRepository {
	return &SwiftCodeRepository{db: db}
}

// scanSwiftCode przetwarza wynik zapytania SQL i wypełnia obiekt models.SwiftCode,
// obsługując kolumnę address (NULL -> "UNKNOWN").
func scanSwiftCode(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.SwiftCode, error) {
	swift := &models.SwiftCode{}
	var address sql.NullString

	err := scanner.Scan(&swift.ID, &swift.SwiftCode, &swift.BankName, &address,
		&swift.CountryISO2, &swift.CountryName, &swift.IsHeadquarter, &swift.HeadquarterID)
	if err != nil {
		return nil, err
	}

	if address.Valid {
		swift.Address = address.String
	} else {
		swift.Address = "UNKNOWN"
	}
	return swift, nil
}

// GetBySwiftCode pobiera kod SWIFT na podstawie wartości swift_code.
func (r *SwiftCodeRepository) GetBySwiftCode(code string) (*models.SwiftCode, error) {
	query := `
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes
		WHERE swift_code = $1;`

	swift, err := scanSwiftCode(r.db.QueryRow(query, code))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Brak wyniku – nie logujemy
		}
		log.Println("Błąd zapytania do bazy danych w GetBySwiftCode:", err)
		return nil, err
	}
	return swift, nil
}

// GetByCountryISO2 pobiera listę kodów SWIFT dla danego kraju (ISO-2).
func (r *SwiftCodeRepository) GetByCountryISO2(countryISO2 string) ([]models.SwiftCode, error) {
	query := `
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes
		WHERE country_iso2 = $1;`

	rows, err := r.db.Query(query, countryISO2)
	if err != nil {
		log.Println("Błąd zapytania do bazy danych w GetByCountryISO2:", err)
		return nil, err
	}
	defer rows.Close()

	var swiftCodes []models.SwiftCode
	for rows.Next() {
		swift, err := scanSwiftCode(rows)
		if err != nil {
			return nil, err
		}
		swiftCodes = append(swiftCodes, *swift)
	}

	if len(swiftCodes) == 0 {
		return nil, sql.ErrNoRows // Brak wyników – nie logujemy
	}
	return swiftCodes, nil
}

func (r *SwiftCodeRepository) DeleteSwiftCode(code string) error {
	query := "DELETE FROM swift_codes WHERE swift_code = $1;"
	if _, err := r.db.Exec(query, code); err != nil {
		log.Println("Błąd usuwania SWIFT code:", err)
		return err
	}
	return nil
}

// DetachBranchesFromHeadquarter odłącza wszystkie branche od danego headquartera.
func (r *SwiftCodeRepository) DetachBranchesFromHeadquarter(headquarterID int64) error {
	query := "UPDATE swift_codes SET headquarter_id = NULL WHERE headquarter_id = $1;"
	if _, err := r.db.Exec(query, headquarterID); err != nil {
		log.Println("Błąd odłączania branchy w DetachBranchesFromHeadquarter:", err)
		return err
	}
	return nil
}

// InsertSwiftCode wstawia nowy rekord kodu SWIFT do bazy danych.
func (r *SwiftCodeRepository) InsertSwiftCode(swift *models.SwiftCode) error {
	query := `
		INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	err := r.db.QueryRow(query, swift.SwiftCode, swift.BankName, swift.Address,
		swift.CountryISO2, swift.CountryName, swift.IsHeadquarter, swift.HeadquarterID).
		Scan(&swift.ID)
	if err != nil {
		log.Println("Błąd wstawiania nowego SWIFT code w InsertSwiftCode:", err) // Logowanie błędów wykonania zapytania
	}
	return err
}

// GetBranchesByHeadquarter pobiera listę branchy powiązanych z danym headquarterem.
func (r *SwiftCodeRepository) GetBranchesByHeadquarter(headquarterCode string) ([]models.SwiftCode, error) {
	var headquarterID int

	// Pobieramy ID headquartera
	err := r.db.QueryRow("SELECT id FROM swift_codes WHERE swift_code = $1;", headquarterCode).Scan(&headquarterID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Brak headquartera – nie logujemy
		}
		log.Println("Błąd pobierania ID headquartera w GetBranchesByHeadquarter:", err) // Logowanie krytycznych błędów
		return nil, err
	}

	// Pobieramy branchy powiązane z headquarterem
	query := `
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes
		WHERE headquarter_id = $1;`
	rows, err := r.db.Query(query, headquarterID)
	if err != nil {
		log.Println("Błąd pobierania branchy w GetBranchesByHeadquarter:", err) // Logowanie błędów
		return nil, err
	}
	defer rows.Close()

	var branches []models.SwiftCode
	for rows.Next() {
		swift, err := scanSwiftCode(rows)
		if err != nil {
			return nil, err // Błędy skanowania zwracamy bez logowania (obsługa w wyższej warstwie)
		}
		branches = append(branches, *swift)
	}

	if len(branches) == 0 {
		return nil, sql.ErrNoRows
	}
	return branches, nil
}
