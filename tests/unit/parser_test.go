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
		{"US", "XYZXYZSS123", "BIC11", "US Branch", "5th Avenue", "Los Angeles", "USA", "America/New_York"},
		{"US", "XYZXYZSSXXX", "BIC11", "US Bank HQ", "Wall Street", "New York", "USA", "America/New_York"},
	}
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)
	assert.Len(t, swiftCodes, 4) // Oczekujemy 4 rekordów w tej samej kolejności

	// Sprawdzanie każdego rekordu w kolejności jego występowania
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
	assert.True(t, swiftCodes[0].IsHeadquarter)
	assert.Nil(t, swiftCodes[0].HeadquarterID)

	assert.Equal(t, "ABCDEFSS001", swiftCodes[1].SwiftCode)
	assert.False(t, swiftCodes[1].IsHeadquarter)
	assert.NotNil(t, swiftCodes[1].HeadquarterID)

	// Uwzględniamy indeksowanie od 1 w bazie (HeadquarterID = ID headquarters + 1)
	assert.Equal(t, swiftCodes[0].ID+1, *swiftCodes[1].HeadquarterID)

	assert.Equal(t, "XYZXYZSS123", swiftCodes[2].SwiftCode)
	assert.False(t, swiftCodes[2].IsHeadquarter)
	assert.NotNil(t, swiftCodes[2].HeadquarterID)
	assert.Equal(t, swiftCodes[3].ID+1, *swiftCodes[2].HeadquarterID)

	assert.Equal(t, "XYZXYZSSXXX", swiftCodes[3].SwiftCode)
	assert.True(t, swiftCodes[3].IsHeadquarter)
	assert.Nil(t, swiftCodes[3].HeadquarterID)
}

func TestParseSwiftCodesInvalidData(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		{"", "AACDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},
		{"PL", "", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},
		{"PL", "BBCDEFSSXXX", "", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},
		{"PL", "CCCDEFSSXXX", "BIC11", "", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},
	}

	testFilePath := "invalid_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)

	// Sprawdzamy, czy nie przetworzył błędnych rekordów
	assert.Len(t, swiftCodes, 0, "Nieprawidłowe rekordy nie powinny zostać przetworzone")
}

func TestParseSwiftCodesEmptyFile(t *testing.T) {
	testFilePath := "empty_test.csv"
	createTestCSV(testFilePath, [][]string{})
	defer os.Remove(testFilePath)

	_, err := services.ParseSwiftCodes(testFilePath)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, io.EOF), "Oczekiwano io.EOF, otrzymano: %v", err)
}

func TestParseSwiftCodesMixedData(t *testing.T) {
	testData := [][]string{
		{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN NAME", "COUNTRY NAME", "TIME ZONE"},
		{"PL", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},       // Poprawny
		{"", "ABCDEFSSXXX", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},         // Brak COUNTRY ISO2 CODE
		{"PL", "XYZXYZSS123", "BIC11", "Branch Bank", "Branch Street 2", "Krakow", "Poland", "Europe/Warsaw"}, // Poprawny
		{"PL", "", "BIC11", "Bank HQ", "Main Street 1", "Warsaw", "Poland", "Europe/Warsaw"},                  // Brak SWIFT CODE
	}

	testFilePath := "mixed_test.csv"
	err := createTestCSV(testFilePath, testData)
	assert.NoError(t, err)
	defer os.Remove(testFilePath)

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.NoError(t, err)

	// Sprawdzamy, czy poprawnie przetworzył tylko poprawne rekordy
	assert.Len(t, swiftCodes, 2, "Oczekiwano dwóch poprawnych rekordów")
	assert.Equal(t, "ABCDEFSSXXX", swiftCodes[0].SwiftCode)
	assert.Equal(t, "XYZXYZSS123", swiftCodes[1].SwiftCode)
}

func TestParseSwiftCodesInvalidFile(t *testing.T) {
	testFilePath := "non_existent.csv"

	swiftCodes, err := services.ParseSwiftCodes(testFilePath)
	assert.Error(t, err)
	assert.Nil(t, swiftCodes)
}
