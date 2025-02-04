package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mroczekDNF/swift-api/internal/repositories"
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
