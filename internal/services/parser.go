package services

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/mroczekDNF/swift-api/internal/models"
)

var idCounter int64 = 1 // Licznik ID, zaczynamy od 1

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

	headquartersMap := make(map[string]models.SwiftCode) // Mapowanie kodu SWIFT na headquarter
	var branches []models.SwiftCode                      // Oddziały (branches)
	var swiftCodes []models.SwiftCode

	for _, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1]))
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))

		swift := models.SwiftCode{
			ID:            idCounter, // Ręcznie przypisujemy unikalne ID
			SwiftCode:     swiftCode,
			BankName:      strings.TrimSpace(record[3]),
			Address:       strings.TrimSpace(record[4]),
			CountryISO2:   countryISO2,
			CountryName:   strings.TrimSpace(record[6]),
			IsHeadquarter: isHeadquarter,
			HeadquarterID: nil, // Jawnie przypisujemy `nil`
		}
		idCounter++ // Zwiększamy licznik ID

		if isHeadquarter {
			// Zapisujemy headquarters w mapie
			swiftCodes = append(swiftCodes, swift)
			headquartersMap[swiftCode[:8]] = swift
		} else {
			// Dodajemy branche do oddzielnej listy
			branches = append(branches, swift)
		}
	}

	// Aktualizujemy branche o poprawne headquarter_id
	for i, branch := range branches {
		if hq, exists := headquartersMap[branch.SwiftCode[:8]]; exists {
			branches[i].HeadquarterID = &hq.ID
		}
	}

	// Dodajemy branche na końcu
	swiftCodes = append(swiftCodes, branches...)

	return swiftCodes, nil
}
