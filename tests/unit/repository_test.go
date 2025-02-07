package unit

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mroczekDNF/swift-api/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// TestGetBySwiftCode - handling a record with a full address
func TestGetBySwiftCode_WithAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
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
	assert.NoError(t, err)
	assert.Equal(t, "123 Bank Street", swift.Address)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetBySwiftCode_NoAddress - handling a record without an address
func TestGetBySwiftCode_NoAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	code := "ABC123XXX"

	query := regexp.QuoteMeta(`
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE swift_code = $1;`)

	rows := sqlmock.NewRows([]string{
		"id", "swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter", "headquarter_id",
	}).AddRow(1, code, "Bank A", nil, "US", "United States", true, nil)

	mock.ExpectQuery(query).WithArgs(code).WillReturnRows(rows)

	swift, err := repo.GetBySwiftCode(code)
	assert.NoError(t, err)
	assert.Equal(t, "UNKNOWN", swift.Address)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByCountryISO2 - handling a list of SWIFT codes with and without addresses
func TestGetByCountryISO2(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	countryISO2 := "US"

	query := regexp.QuoteMeta(`
		SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter, headquarter_id
		FROM swift_codes WHERE country_iso2 = $1;`)

	rows := sqlmock.NewRows([]string{
		"id", "swift_code", "bank_name", "address", "country_iso2", "country_name", "is_headquarter", "headquarter_id",
	}).
		AddRow(1, "ABC123XXX", "Bank A", "123 Bank Street", "US", "United States", true, nil).
		AddRow(2, "DEF456", "Bank B", nil, "US", "United States", false, 1)

	mock.ExpectQuery(query).WithArgs(countryISO2).WillReturnRows(rows)

	swiftCodes, err := repo.GetByCountryISO2(countryISO2)
	assert.NoError(t, err)
	assert.Len(t, swiftCodes, 2)
	assert.Equal(t, "123 Bank Street", swiftCodes[0].Address)
	assert.Equal(t, "UNKNOWN", swiftCodes[1].Address)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDetachBranchesFromHeadquarter - detaching branches from a headquarter
func TestDetachBranchesFromHeadquarter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewSwiftCodeRepository(db)
	headquarterID := int64(1)

	query := regexp.QuoteMeta("UPDATE swift_codes SET headquarter_id = NULL WHERE headquarter_id = $1;")
	mock.ExpectExec(query).WithArgs(headquarterID).WillReturnResult(sqlmock.NewResult(0, 2))

	err = repo.DetachBranchesFromHeadquarter(headquarterID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
