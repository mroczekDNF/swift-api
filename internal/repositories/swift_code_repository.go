package repositories

import (
	"database/sql"
	"log"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// SwiftCodeRepositoryInterface defines repository methods
type SwiftCodeRepositoryInterface interface {
	GetBySwiftCode(code string) (*models.SwiftCode, error)
	GetByCountryISO2(countryISO2 string) ([]models.SwiftCode, error)
	DeleteSwiftCode(code string) error
	DetachBranchesFromHeadquarter(headquarterID int64) error
	InsertSwiftCode(swift *models.SwiftCode) error
	GetBranchesByHeadquarter(headquarterCode string) ([]models.SwiftCode, error)
}

// SwiftCodeRepository handles operations on the swift_codes table
type SwiftCodeRepository struct {
	db *sql.DB
}

// NewSwiftCodeRepository creates a new SwiftCode repository instance
func NewSwiftCodeRepository(db *sql.DB) *SwiftCodeRepository {
	return &SwiftCodeRepository{db: db}
}

// scanSwiftCode processes the SQL query result and populates a models.SwiftCode object
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

	swift.Address = "UNKNOWN"
	if address.Valid {
		swift.Address = address.String
	}
	return swift, nil
}

// GetBySwiftCode retrieves a SWIFT code by its value
func (r *SwiftCodeRepository) GetBySwiftCode(code string) (*models.SwiftCode, error) {
	query := "SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id FROM swift_codes WHERE swift_code = $1;"

	swift, err := scanSwiftCode(r.db.QueryRow(query, code))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Database query error in GetBySwiftCode:", err)
	}
	return swift, err
}

// GetByCountryISO2 retrieves a list of SWIFT codes for a given country
func (r *SwiftCodeRepository) GetByCountryISO2(countryISO2 string) ([]models.SwiftCode, error) {
	query := "SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id FROM swift_codes WHERE country_iso2 = $1;"

	rows, err := r.db.Query(query, countryISO2)
	if err != nil {
		log.Println("Database query error in GetByCountryISO2:", err)
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
		return nil, sql.ErrNoRows
	}
	return swiftCodes, nil
}

// DeleteSwiftCode removes a SWIFT code from the database
func (r *SwiftCodeRepository) DeleteSwiftCode(code string) error {
	query := "DELETE FROM swift_codes WHERE swift_code = $1;"
	_, err := r.db.Exec(query, code)
	if err != nil {
		log.Println("Error deleting SWIFT code:", err)
	}
	return err
}

// DetachBranchesFromHeadquarter detaches all branches from a given headquarter
func (r *SwiftCodeRepository) DetachBranchesFromHeadquarter(headquarterID int64) error {
	query := "UPDATE swift_codes SET headquarter_id = NULL WHERE headquarter_id = $1;"
	_, err := r.db.Exec(query, headquarterID)
	if err != nil {
		log.Println("Error detaching branches in DetachBranchesFromHeadquarter:", err)
	}
	return err
}

// InsertSwiftCode inserts a new SWIFT code record into the database
func (r *SwiftCodeRepository) InsertSwiftCode(swift *models.SwiftCode) error {
	query := "INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"

	err := r.db.QueryRow(query, swift.SwiftCode, swift.BankName, swift.Address,
		swift.CountryISO2, swift.CountryName, swift.IsHeadquarter, swift.HeadquarterID).Scan(&swift.ID)
	if err != nil {
		log.Println("Error inserting new SWIFT code in InsertSwiftCode:", err)
	}
	return err
}

// GetBranchesByHeadquarter retrieves branches associated with a given headquarter
func (r *SwiftCodeRepository) GetBranchesByHeadquarter(headquarterCode string) ([]models.SwiftCode, error) {
	var headquarterID int

	// Retrieve headquarter ID
	err := r.db.QueryRow("SELECT id FROM swift_codes WHERE swift_code = $1;", headquarterCode).Scan(&headquarterID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching headquarter ID in GetBranchesByHeadquarter:", err)
		return nil, err
	}

	// Retrieve branches associated with the headquarter
	query := "SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id FROM swift_codes WHERE headquarter_id = $1;"
	rows, err := r.db.Query(query, headquarterID)
	if err != nil {
		log.Println("Error fetching branches in GetBranchesByHeadquarter:", err)
		return nil, err
	}
	defer rows.Close()

	var branches []models.SwiftCode
	for rows.Next() {
		swift, err := scanSwiftCode(rows)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *swift)
	}

	if len(branches) == 0 {
		return nil, sql.ErrNoRows
	}
	return branches, nil
}
