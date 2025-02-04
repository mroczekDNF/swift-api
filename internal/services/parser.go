package services

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mroczekDNF/swift-api/internal/models"
)

const (
	ColumnCountryISO2 = 0
	ColumnSwiftCode   = 1
	ColumnBankName    = 3
	ColumnAddress     = 4
	ColumnCountryName = 6
	CodeType          = 2
)

// isValidRecord checks if the record contains the required data.
func isValidRecord(record []string) bool {
	if len(record) < 7 {
		log.Println("Rejected record: insufficient data")
		return false
	}

	swiftCode := strings.TrimSpace(record[ColumnSwiftCode])
	countryISO2 := strings.TrimSpace(record[ColumnCountryISO2])
	bankName := strings.TrimSpace(record[ColumnBankName])
	codeType := strings.TrimSpace(record[CodeType])

	if codeType == "" {
		log.Println("Rejected record:", swiftCode, "- missing SWIFT code type")
		return false
	}
	if len(countryISO2) != 2 {
		log.Println("Rejected record:", swiftCode, "- invalid country code:", countryISO2)
		return false
	}
	if len(swiftCode) < 8 || len(swiftCode) > 11 {
		log.Println("Rejected record:", swiftCode, "- invalid SWIFT code")
		return false
	}
	if bankName == "" {
		log.Println("Rejected record:", swiftCode, "- missing bank name")
		return false
	}
	return true
}

// filterValidRecords filters valid records from raw data.
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
// Uses errors.Is to check for EOF and returns an error if another issue occurs.
func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Read headers
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
		if strings.HasSuffix(swiftCode, "XXX") && len(swiftCode) >= 8 {
			key := swiftCode[:8]
			if _, exists := headquartersMap[key]; !exists {
				headquartersMap[key] = idCounter
				idCounter++
			}
		}
	}
	return headquartersMap
}

// processValidRecords processes valid records, assigning them IDs and relationships between branches and headquarters.
// Since IDs are 1-based (starting from 1), records are stored in a preallocated slice at position (ID - 1).
func processValidRecords(validData [][]string, headquartersMap map[string]int64) []models.SwiftCode {
	// Preallocate slice based on the number of records.
	swiftCodes := make([]models.SwiftCode, len(validData))

	// If there are already 3 headquarters, the first branch should get ID = 4.
	var idBranch int64 = int64(len(headquartersMap)) + 1

	for _, record := range validData {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[ColumnSwiftCode]))
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[ColumnCountryISO2]))

		swift := models.SwiftCode{
			ID:            -1, // Temporary value, will be replaced
			SwiftCode:     swiftCode,
			BankName:      strings.TrimSpace(record[ColumnBankName]),
			Address:       strings.TrimSpace(record[ColumnAddress]),
			CountryISO2:   countryISO2,
			CountryName:   strings.TrimSpace(record[ColumnCountryName]),
			IsHeadquarter: isHeadquarter,
			HeadquarterID: nil,
		}

		if isHeadquarter {
			key := swiftCode[:8]
			if hqID, exists := headquartersMap[key]; exists {
				swift.ID = hqID
			}
		} else {
			// Assign HeadquarterID if a matching headquarters exists.
			key := swiftCode[:8]
			if hqID, exists := headquartersMap[key]; exists {
				swift.HeadquarterID = &hqID
			} else {
				// If a branch does not have a matching headquarters, log a warning.
				log.Printf("Warning: Branch %s does not have a corresponding headquarters in the map.", swiftCode)
			}
			swift.ID = idBranch
			idBranch++
		}

		// Store the record in the preallocated slice at position (ID - 1) since IDs are 1-based.
		swiftCodes[swift.ID-1] = swift
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
