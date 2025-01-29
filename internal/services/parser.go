package services

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// ParseSwiftCodes parsuje dane SWIFT z pliku CSV i zwraca listę struktur SwiftCode
func ParseSwiftCodes(filePath string) ([]models.SwiftCode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(data) <= 1 {
		return nil, nil
	}

	data = data[1:]

	headquartersMap := make(map[string]int64)
	var swiftCodes []models.SwiftCode

	for _, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1]))
		if strings.HasSuffix(swiftCode, "XXX") {
			headquartersMap[swiftCode[:8]] = 0
		}
	}

	for _, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1]))
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))

		var headquarterID *int64
		if !isHeadquarter {
			if id, exists := headquartersMap[swiftCode[:8]]; exists {
				headquarterID = &id
			}
		}

		swift := models.SwiftCode{
			SwiftCode:     swiftCode,
			BankName:      strings.TrimSpace(record[3]),
			Address:       strings.TrimSpace(record[4]),
			CountryISO2:   countryISO2,
			CountryName:   strings.TrimSpace(record[6]),
			IsHeadquarter: isHeadquarter,
			HeadquarterID: headquarterID,
		}
		swiftCodes = append(swiftCodes, swift)
	}
	return swiftCodes, nil
}
