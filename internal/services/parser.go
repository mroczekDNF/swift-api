package services

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// ParseSwiftCodes parsuje dane SWIFT z pliku CSV i zwraca listę struktur SwiftCode
func ParseSwiftCodes(filePath string) ([]models.SwiftCode, error) {
	// Otwórz plik CSV
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Wczytaj dane z pliku CSV
	reader := csv.NewReader(file)
	reader.Comma = ',' // Zmień separator, jeśli plik używa czegoś innego
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Sprawdź, czy są dane w pliku
	if len(data) <= 1 {
		return nil, nil // Plik jest pusty lub zawiera tylko nagłówki
	}

	// Pomijamy nagłówki (pierwszy wiersz)
	data = data[1:]

	// Mapa headquarters
	headquartersMap := make(map[string]string) // Klucz: pierwsze 8 znaków SWIFT, Wartość: pełny SWIFT code headquarter

	// Lista wszystkich SWIFT codes
	var swiftCodes []models.SwiftCode

	// Pierwszy przebieg: znajdź wszystkie headquarters
	for _, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1])) // Kolumna B: SWIFT CODE

		// Sprawdź, czy to headquarter
		if !strings.HasSuffix(swiftCode, "XXX") {
			continue // Jeśli nie headquarter, pomiń
		}

		// Dodaj headquarter do mapy
		headquartersMap[swiftCode[:8]] = swiftCode
	}

	// Drugi przebieg: przetwarzanie wszystkich rekordów i przypisanie headquarter_id
	for _, record := range data {
		swiftCode := strings.ToUpper(strings.TrimSpace(record[1])) // Kolumna B: SWIFT CODE
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

		// Walidacja kodu ISO-2
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))
		if len(countryISO2) != 2 {
			log.Printf("Invalid country ISO2 code: %s, skipping record", countryISO2)
			continue
		}

		// Znajdź headquarter_id dla branchy
		var headquarterID *string
		if !isHeadquarter {
			if hqSwift, exists := headquartersMap[swiftCode[:8]]; exists {
				headquarterID = &hqSwift
			}
		}

		// Tworzenie struktury SwiftCode
		swift := models.SwiftCode{
			SwiftCode:     swiftCode,
			BankName:      strings.ToUpper(strings.TrimSpace(record[3])), // Kolumna D: NAME
			Address:       strings.ToUpper(strings.TrimSpace(record[4])), // Kolumna E: ADDRESS
			CountryISO2:   countryISO2,
			CountryName:   strings.ToUpper(strings.TrimSpace(record[6])), // Kolumna G: COUNTRY NAME
			IsHeadquarter: isHeadquarter,
			HeadquarterID: headquarterID,
		}

		// Dodaj rekord do listy
		swiftCodes = append(swiftCodes, swift)
	}
	return swiftCodes, nil
}
