package services

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mroczekDNF/swift-api/internal/models"
)

const (
	ColumnCountryISO2  = 0
	ColumnSwiftCode    = 1
	ColumnBankName     = 3
	ColumnAddress      = 4
	ColumnCountryName  = 6
	CodeType           = 2
	headquartersSuffix = "XXX"
)

// Precompiled regular expressions
var (
	// SWIFT: 4 letters, 2 letters, 2 alphanumeric characters, and optionally 3 alphanumeric characters.
	swiftCodeRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	// Country ISO2: exactly 2 uppercase letters.
	countryISO2Regex = regexp.MustCompile(`^[A-Z]{2}$`)
)

// isValidRecord checks if the record contains the required data and validates its format.
func isValidRecord(record []string) bool {
	if len(record) < 7 {
		log.Println("Rejected record: insufficient data")
		return false
	}

	// Normalize data
	swiftCode := strings.ToUpper(strings.TrimSpace(record[ColumnSwiftCode]))
	countryISO2 := strings.ToUpper(strings.TrimSpace(record[ColumnCountryISO2]))
	bankName := strings.TrimSpace(record[ColumnBankName])
	codeType := strings.TrimSpace(record[CodeType])

	if codeType == "" {
		log.Println("Rejected record:", swiftCode, "- missing SWIFT code type")
		return false
	}

	// Validate the country code using regex.
	if !countryISO2Regex.MatchString(countryISO2) {
		log.Println("Rejected record:", swiftCode, "- invalid country code:", countryISO2)
		return false
	}

	if len(swiftCode) < 8 || len(swiftCode) > 11 {
		log.Println("Rejected record:", swiftCode, "- invalid SWIFT code length")
		return false
	}
	// Validate the SWIFT code format using regex.
	if !swiftCodeRegex.MatchString(swiftCode) {
		log.Println("Rejected record:", swiftCode, "- invalid SWIFT code format")
		return false
	}

	if bankName == "" {
		log.Println("Rejected record:", swiftCode, "- missing bank name")
		return false
	}
	return true
}

// filterValidRecords filters valid records from the raw data.
func filterValidRecords(data [][]string) [][]string {
	validRecords := make([][]string, 0, len(data))
	for _, record := range data {
		if isValidRecord(record) {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

// readCSV reads data from a CSV file and returns a list of records.
func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Read headers.
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return [][]string{}, nil
	}

	var records [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// createHeadquartersMap creates a map of headquarters, assigning them unique IDs starting from 1.
// The key in the map is the first 8 characters of the SWIFT code.
func createHeadquartersMap(validData [][]string) map[string]int64 {
	headquartersMap := make(map[string]int64)
	var idCounter int64 = 1 // IDs start from 1

	for _, record := range validData {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[ColumnSwiftCode]))
		// A record is considered a headquarters if the SWIFT code ends with "XXX" and has at least 8 characters.
		if strings.HasSuffix(swiftCode, headquartersSuffix) && len(swiftCode) >= 8 {
			key := swiftCode[:8]
			if _, exists := headquartersMap[key]; !exists {
				headquartersMap[key] = idCounter
			}
		}
		idCounter++
	}
	return headquartersMap
}

// processValidRecords processes valid records, assigning them IDs and establishing relationships
// between branches and headquarters.
func processValidRecords(validData [][]string, headquartersMap map[string]int64) []models.SwiftCode {
	// Preallocate slice based on the number of records.
	swiftCodes := make([]models.SwiftCode, len(validData))

	for idx, record := range validData {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[ColumnSwiftCode]))
		isHeadquarter := strings.HasSuffix(swiftCode, headquartersSuffix)
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[ColumnCountryISO2]))

		// Sprawdź, czy adres jest pusty, i ustaw "UNKNOWN" w takim przypadku
		address := strings.TrimSpace(record[ColumnAddress])
		if address == "" {
			address = "UNKNOWN"
		}

		swift := models.SwiftCode{
			ID:            int64(idx), // Temporary value, which may be replaced later.
			SwiftCode:     swiftCode,
			BankName:      strings.TrimSpace(record[ColumnBankName]),
			Address:       address,
			CountryISO2:   countryISO2,
			CountryName:   strings.TrimSpace(record[ColumnCountryName]),
			IsHeadquarter: isHeadquarter,
			HeadquarterID: nil,
		}

		if !isHeadquarter {
			key := swiftCode[:8]
			if hqID, exists := headquartersMap[key]; exists {
				swift.HeadquarterID = &hqID
			}
		}
		swiftCodes[idx] = swift
	}

	return swiftCodes
}

// ParseSwiftCodes is the main function for parsing SWIFT data from a CSV file.
func ParseSwiftCodes(filePath string) ([]models.SwiftCode, error) {
	data, err := readCSV(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.SwiftCode{}, nil
	}

	validData := filterValidRecords(data)
	if len(validData) == 0 {
		return []models.SwiftCode{}, nil
	}

	headquartersMap := createHeadquartersMap(validData)
	swiftCodes := processValidRecords(validData, headquartersMap)

	return swiftCodes, nil
}
