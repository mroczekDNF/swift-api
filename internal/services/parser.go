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

	// Mapa do przechowywania ID dla headquarters
	headquartersMap := make(map[string]int64)

	// Prealokacja pamięci dla listy wynikowej (eliminacja `append()`)
	swiftCodes := make([]models.SwiftCode, len(data))

	// 🔹 **Pierwsza pętla** – przypisujemy ID tylko dla headquarters
	for _, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1]))

		if strings.HasSuffix(swiftCode, "XXX") {
			headquartersMap[swiftCode[:8]] = idCounter // Zapisujemy ID dla headquarters
			idCounter++                                // Inkrementujemy ID tylko dla headquarters
		}
	}

	// 🔹 **Druga pętla** – teraz tworzymy pełne instancje `SwiftCode`
	for i, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1]))
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))

		// Tworzymy instancję `SwiftCode`
		swift := models.SwiftCode{
			ID:            idCounter,
			SwiftCode:     swiftCode,
			BankName:      strings.TrimSpace(record[3]),
			Address:       strings.TrimSpace(record[4]),
			CountryISO2:   countryISO2,
			CountryName:   strings.TrimSpace(record[6]),
			IsHeadquarter: isHeadquarter,
			HeadquarterID: nil, // Jawnie przypisujemy `nil`
		}

		if isHeadquarter {
			// Headquarters używają ID z mapy
			swift.ID = headquartersMap[swiftCode[:8]]
		} else {
			// Branch używa nowego ID i przypisuje ID headquarters
			if hqID, exists := headquartersMap[swiftCode[:8]]; exists {
				swift.HeadquarterID = &hqID
			}
			idCounter++ // Inkrementujemy ID tylko dla branchy
		}

		// Zapisujemy rekord w prealokowanej tablicy
		swiftCodes[i] = swift
	}

	return swiftCodes, nil
}
