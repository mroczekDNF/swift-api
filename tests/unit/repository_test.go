package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// TestGetBySwiftCode - obsługa rekordu z pełnym adresem
func TestGetBySwiftCode_WithAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Błąd przy tworzeniu sqlmock: %v", err)
	}
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	code := "ABC123XXX"

	query := regexp.QuoteMeta(`
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE swift_code = $1;`)

	rows := sqlmock.NewRows([]string{
		"id", "swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter", "headquarter_id",
	}).AddRow(1, code, "Bank A", "123 Bank Street", "US", "United States", true, nil)

	mock.ExpectQuery(query).WithArgs(code).WillReturnRows(rows)

	swift, err := repo.GetBySwiftCode(code)
	if err != nil {
		t.Errorf("Nieoczekiwany błąd: %v", err)
	}
	if swift.Address != "123 Bank Street" {
		t.Errorf("Oczekiwano adresu '123 Bank Street', otrzymano: %s", swift.Address)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Niespełnione oczekiwania: %v", err)
	}
}

// TestGetBySwiftCode_NoAddress - obsługa rekordu bez adresu
func TestGetBySwiftCode_NoAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Błąd przy tworzeniu sqlmock: %v", err)
	}
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	code := "ABC123XXX"

	query := regexp.QuoteMeta(`
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE swift_code = $1;`)

	// Przekazujemy NULL jako adres
	rows := sqlmock.NewRows([]string{
		"id", "swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter", "headquarter_id",
	}).AddRow(1, code, "Bank A", nil, "US", "United States", true, nil) // <- NULL jako address

	mock.ExpectQuery(query).WithArgs(code).WillReturnRows(rows)

	swift, err := repo.GetBySwiftCode(code)
	if err != nil {
		t.Errorf("Nieoczekiwany błąd: %v", err)
	}

	// Sprawdzamy, czy brak adresu w bazie zwraca poprawnie "UNKNOWN"
	if swift.Address != "UNKNOWN" {
		t.Errorf("Oczekiwano 'UNKNOWN' dla pustego adresu, otrzymano: %s", swift.Address)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Niespełnione oczekiwania: %v", err)
	}
}

// TestGetByCountryISO2 - obsługa listy SWIFT codes z adresami i bez adresów
func TestGetByCountryISO2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Błąd przy tworzeniu sqlmock: %v", err)
	}
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	countryISO2 := "US"

	query := regexp.QuoteMeta(`
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE country_iso2 = $1;`)

	rows := sqlmock.NewRows([]string{
		"id", "swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter", "headquarter_id",
	}).
		// Rekord z adresem
		AddRow(1, "ABC123XXX", "Bank A", "123 Bank Street", "US", "United States", true, nil).
		// Rekord bez adresu (NULL w bazie)
		AddRow(2, "DEF456", "Bank B", nil, "US", "United States", false, 1)

	mock.ExpectQuery(query).WithArgs(countryISO2).WillReturnRows(rows)

	swiftCodes, err := repo.GetByCountryISO2(countryISO2)
	if err != nil {
		t.Errorf("Nieoczekiwany błąd: %v", err)
	}
	if len(swiftCodes) != 2 {
		t.Errorf("Oczekiwano 2 rekordów, otrzymano: %d", len(swiftCodes))
	}

	// Sprawdzamy, czy pierwszy rekord ma właściwy adres
	if swiftCodes[0].Address != "123 Bank Street" {
		t.Errorf("Oczekiwano '123 Bank Street', otrzymano: %s", swiftCodes[0].Address)
	}

	// Sprawdzamy, czy drugi rekord bez adresu zwraca `"UNKNOWN"`
	if swiftCodes[1].Address != "UNKNOWN" {
		t.Errorf("Oczekiwano 'UNKNOWN', otrzymano: %s", swiftCodes[1].Address)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Niespełnione oczekiwania: %v", err)
	}
}

func TestDetachBranchesFromHeadquarter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	headquarterID := int64(1)

	// Mockujemy zapytanie
	query := regexp.QuoteMeta("UPDATE swift_codes SET headquarter_id = NULL WHERE headquarter_id = $1;")
	mock.ExpectExec(query).WithArgs(headquarterID).WillReturnResult(sqlmock.NewResult(0, 2)) // 2 rekordy zostały zaktualizowane

	// Test funkcji
	err = repo.DetachBranchesFromHeadquarter(headquarterID)
	assert.NoError(t, err)

	// Weryfikacja oczekiwań
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertSwiftCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	swift := &models.SwiftCode{
		SwiftCode:     "ABCDEFSSXXX",
		BankName:      "Bank HQ",
		Address:       "Main Street 1",
		CountryISO2:   "PL",
		CountryName:   "Poland",
		IsHeadquarter: true,
	}

	// Mockujemy zapytanie
	query := regexp.QuoteMeta(`INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`)
	mock.ExpectQuery(query).
		WithArgs(swift.SwiftCode, swift.BankName, swift.Address, swift.CountryISO2, swift.CountryName, swift.IsHeadquarter, swift.HeadquarterID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Zwracamy ID = 1

	// Test funkcji
	err = repo.InsertSwiftCode(swift)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), swift.ID)

	// Weryfikacja oczekiwań
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBranchesByHeadquarter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	headquarterCode := "ABCDEFSSXXX"
	headquarterID := 1

	// Mock zapytania do pobrania ID headquarters
	headquarterQuery := regexp.QuoteMeta("SELECT id FROM swift_codes WHERE swift_code = $1;")
	mock.ExpectQuery(headquarterQuery).WithArgs(headquarterCode).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(headquarterID))

	// Mock zapytania do pobrania branchy
	branchesQuery := regexp.QuoteMeta(`
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE headquarter_id = $1;`)
	mock.ExpectQuery(branchesQuery).WithArgs(headquarterID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter", "headquarter_id",
		}).
			AddRow(2, "ABCDEFSS001", "Branch Bank", "Branch Street 1", "PL", "Poland", false, headquarterID).
			AddRow(3, "ABCDEFSS002", "Branch Bank 2", nil, "PL", "Poland", false, headquarterID)) // Adres = NULL

	// Test funkcji
	branches, err := repo.GetBranchesByHeadquarter(headquarterCode)
	assert.NoError(t, err)
	assert.Len(t, branches, 2)

	// Sprawdzanie pierwszego rekordu
	assert.Equal(t, int64(2), branches[0].ID)
	assert.Equal(t, "ABCDEFSS001", branches[0].SwiftCode)
	assert.Equal(t, "Branch Street 1", branches[0].Address)

	// Sprawdzanie drugiego rekordu (adres = UNKNOWN)
	assert.Equal(t, int64(3), branches[1].ID)
	assert.Equal(t, "ABCDEFSS002", branches[1].SwiftCode)
	assert.Equal(t, "UNKNOWN", branches[1].Address)

	// Weryfikacja oczekiwań
	assert.NoError(t, mock.ExpectationsWereMet())
}
