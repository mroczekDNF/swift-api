package unit

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/mroczekDNF/swift-api/internal/services"
	"github.com/stretchr/testify/assert"
)

// createTestCSV creates a temporary CSV file with the given data.
func createTestCSV(filePath string, data [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(data)
}

// TestParseSwiftCodes tests that valid records are processed correctly.
func TestParseSwiftCodes(t *testing.T) {
	testFilePath := "test_swift_codes.csv"
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		{"PL", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},
		{"PL", "ABCDEFSS001", "BIC11", "Branch Bank", "Branch Street 2", "Krakow", "Poland", "Europe/Warsaw"},
		{"US", "XYZXYZSS123", "BIC11", "US Branch", "5th Avenue", "Los Angeles", "USA", "America/New_York"},
		{"US", "XYZXYZSSXXX", "BIC11", "US Bank HQ", "Wall Street", "New York", "USA", "America/New_York"},
	}
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	// We expect 4 records (header not counted) in the same order.
	assert.Len(t, swiftCodes, 4)

	// Checking each record in the order they are processed.
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
	assert.True(t, swiftCodes[0].IsHeadquarter)
	assert.Nil(t, swiftCodes[0].HeadquarterID)

	assert.Equal(t, "ABCDEFSS001", swiftCodes[1].SwiftCode)
	assert.False(t, swiftCodes[1].IsHeadquarter)
	assert.NotNil(t, swiftCodes[1].HeadquarterID)
	// The HeadquarterID should equal the headquarters ID (the first record's ID plus 1).
	assert.Equal(t, swiftCodes[0].ID+1, *swiftCodes[1].HeadquarterID)

	assert.Equal(t, "XYZXYZSS123", swiftCodes[2].SwiftCode)
	assert.False(t, swiftCodes[2].IsHeadquarter)
	assert.NotNil(t, swiftCodes[2].HeadquarterID)
	assert.Equal(t, swiftCodes[3].ID+1, *swiftCodes[2].HeadquarterID)

	assert.Equal(t, "XYZXYZSSXXX", swiftCodes[3].SwiftCode)
	assert.True(t, swiftCodes[3].IsHeadquarter)
	assert.Nil(t, swiftCodes[3].HeadquarterID)
}

// TestParseSwiftCodesInvalidData tests that records with missing required fields are rejected.
func TestParseSwiftCodesInvalidData(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		{"", "AACDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"}, // Missing country code.
		{"PL", "", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},          // Missing SWIFT code.
		{"PL", "BBCDEFSSXXX", "", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},    // Missing code type.
		{"PL", "CCCDEFSSXXX", "BIC11", "", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},      // Missing bank name.
	}

	testFilePath := "invalid_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	// We expect no records to be processed if all records are invalid.
	assert.Len(t, swiftCodes, 0, "Invalid records should not be processed")
}

// TestParseSwiftCodesEmptyFile tests that parsing an empty CSV file returns an error.
func TestParseSwiftCodesEmptyFile(t *testing.T) {
	testFilePath := "empty_test.csv"
	createTestCSV(testFilePath, [][]string{})
	defer os.Remove(testFilePath)

	_, err := services.ParseSwiftCodes(testFilePath)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF), "Expected io.EOF, got: %v", err)
}

// TestParseSwiftCodesMixedData tests that only valid records are processed from mixed data.
func TestParseSwiftCodesMixedData(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		{"PL", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},       // Valid.
		{"", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},         // Missing country code.
		{"PL", "XYZXYZSS123", "BIC11", "Branch Bank", "Branch Street 2", "Krakow", "Poland", "Europe/Warsaw"}, // Valid.
		{"PL", "", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                  // Missing SWIFT code.
	}

	testFilePath := "mixed_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	// Only the two valid records should be processed.
	assert.Len(t, swiftCodes, 2, "Expected two valid records")
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
	assert.Equal(t, "XYZXYZSS123", swiftCodes[1].SwiftCode)
}

// TestParseSwiftCodesInvalidFile tests that parsing a non-existent file returns an error.
func TestParseSwiftCodesInvalidFile(t *testing.T) {
	testFilePath := "non_existent.csv"

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.Error(t, err)
	assert.Nil(t, swiftCodes)
}

// Additional tests with invalid data.

// Modified test: TestParseSwiftCodesInvalidSwiftFormat tests that records with clearly invalid SWIFT code formats are rejected,
// while records with lowercase letters are normalized and accepted.
func TestParseSwiftCodesInvalidSwiftFormat(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		// Invalid SWIFT format: not matching expected pattern.
		{"PL", "INVALID", "BIC11", "Invalid SWIFT", "Some Address", "City", "Poland", "Europe/Warsaw"},
		// Valid SWIFT format but with lowercase letters; after normalization it should be accepted.
		{"PL", "abcdEFSSXXX", "BIC11", "Valid SWIFT", "Some Address", "City", "Poland", "Europe/Warsaw"},
	}
	testFilePath := "invalid_swift_format.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	// Only the record with valid (normalized) SWIFT code should be processed.
	assert.Len(t, swiftCodes, 1, "Only records with valid normalized SWIFT codes should be processed")
	// Check that the SWIFT code is normalized to uppercase.
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
}

// TestParseSwiftCodesInvalidCountryISO2 tests that records with invalid country ISO2 codes are rejected.
func TestParseSwiftCodesInvalidCountryISO2(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		// Invalid country code: only one letter.
		{"P", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Some Address", "City", "Poland", "Europe/Warsaw"},
		// Invalid country code: three characters.
		{"USA", "XYZXYZSS123", "BIC11", "Branch", "Some Address", "City", "USA", "America/New_York"},
	}
	testFilePath := "invalid_country_iso2.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	// Both records should be rejected.
	assert.Len(t, swiftCodes, 0, "Records with invalid country ISO2 should not be processed")
}

// TestParseSwiftCodesMissingFields tests that records missing required fields are rejected.
func TestParseSwiftCodesMissingFields(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		// Missing NAME field (bank name is empty).
		{"PL", "ABCDEFSSXXX", "BIC11", "", "Some Address", "City", "Poland", "Europe/Warsaw"},
		// Missing CODE TYPE field.
		{"PL", "XYZXYZSS123", "", "Bank", "Some Address", "City", "Poland", "Europe/Warsaw"},
	}
	testFilePath := "missing_fields.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	// Both records should be rejected.
	assert.Len(t, swiftCodes, 0, "Records with missing required fields should not be processed")
}

// TestParseSwiftCodesNormalization tests that the parser normalizes SWIFT codes to uppercase.
func TestParseSwiftCodesNormalization(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		// SWIFT code in mixed case should be normalized.
		{"PL", "aBcDeFssXxX", "BIC11", "Bank HQ", "Some Address", "City", "Poland", "Europe/Warsaw"},
	}
	testFilePath := "normalization_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	assert.Len(t, swiftCodes, 1)
	// Check that the SWIFT code is normalized to uppercase.
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
}
