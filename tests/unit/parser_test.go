package unit

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/internal/services"
	"github.com/stretchr/testify/assert"
)

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

func TestParseSwiftCodes(t *testing.T) {
	testFilePath := "test_swift_codes.csv"
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		{"PL", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},
		{"PL", "ABCDEFSS001", "BIC11", "Branch Bank", "Branch Street 2", "Krakow", "Poland", "Europe/Warsaw"},
		{"US", "XYZXYZSSXXX", "BIC11", "US Bank HQ", "Wall Street", "New York", "USA", "America/New_York"},
		{"US", "XYZXYZSS123", "BIC11", "US Branch", "5th Avenue", "Los Angeles", "USA", "America/New_York"},
	}
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	assert.Len(t, swiftCodes, 4) // Expecting 4 records

	// Validate headquarters
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
	assert.True(t, swiftCodes[0].IsHeadquarter)

	// Ensure HeadquarterID is nil for all headquarters
	for _, code := range swiftCodes {
		if code.IsHeadquarter {
			assert.Nil(t, code.HeadquarterID, "Headquarter %s should have HeadquarterID == nil", code.SwiftCode)
		}
	}

	// Validate branch
	assert.Equal(t, "ABCDEFSS001", swiftCodes[2].SwiftCode)
	assert.False(t, swiftCodes[2].IsHeadquarter)

	// Ensure branch has correct HeadquarterID
	assert.NotNil(t, swiftCodes[2].HeadquarterID, "Branch ABCDEFSS001 should have a HeadquarterID assigned")
	assert.Equal(t, swiftCodes[0].ID, *swiftCodes[2].HeadquarterID, "Branch ABCDEFSS001 should be associated with headquarter ABCDEFSSXXX")

	// Validate branch → headquarter mapping
	branchSwiftCode := "ABCDEFSS001"
	headquarterSwiftCode := "ABCDEFSSXXX"
	var branch, headquarter *models.SwiftCode
	for _, code := range swiftCodes {
		if code.SwiftCode == branchSwiftCode {
			branch = &code
		}
		if code.SwiftCode == headquarterSwiftCode {
			headquarter = &code
		}
	}

	assert.NotNil(t, branch, "Branch %s should exist in the results", branchSwiftCode)
	assert.NotNil(t, headquarter, "Headquarter %s should exist in the results", headquarterSwiftCode)

	assert.Equal(t, headquarter.ID, *branch.HeadquarterID, "Branch %s should have HeadquarterID %d, matching headquarter %s", branchSwiftCode, headquarter.ID, headquarterSwiftCode)
}

func TestParseSwiftCodesInvalidData(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"}, // Headers
		{"", "AACDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                 // Missing COUNTRY ISO2 CODE
		{"PL", "", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                          // Missing SWIFT CODE
		{"PL", "BBCDEFSSXXX", "", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                    // Missing CODE TYPE
		{"PL", "CCCDEFSSXXX", "BIC11", "", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                      // Missing NAME
	}

	testFilePath := "invalid_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)

	// Ensure no invalid records are processed
	assert.Len(t, swiftCodes, 0, "Invalid records should not be processed")
}

func TestParseSwiftCodesEmptyFile(t *testing.T) {
	testFilePath := "empty_test.csv"
	createTestCSV(testFilePath, [][]string{})
	defer os.Remove(testFilePath)

	_, err := services.ParseSwiftCodes(testFilePath)
	assert.Error(t, err)                                                    // Expect an error
	assert.True(t, errors.Is(err, io.EOF), "Expected io.EOF, got: %v", err) // Ensure it's an EOF error
}

func TestParseSwiftCodesMixedData(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"}, // Headers
		{"PL", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},               // Valid
		{"", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                 // Missing COUNTRY ISO2 CODE
		{"PL", "XYZXYZSS123", "BIC11", "Branch Bank", "Branch Street 2", "Krakow", "Poland", "Europe/Warsaw"},         // Valid
		{"PL", "", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                          // Missing SWIFT CODE
	}

	testFilePath := "mixed_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)

	// Ensure only valid records are processed
	assert.Len(t, swiftCodes, 2, "Expected two valid records")
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
	assert.Equal(t, "XYZXYZSS123", swiftCodes[1].SwiftCode)
}

func TestParseSwiftCodesInvalidFile(t *testing.T) {
	testFilePath := "non_existent.csv"

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.Error(t, err)
	assert.Nil(t, swiftCodes)
}
