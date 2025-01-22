package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/mroczekDNF/swift-api/internal/models"
)

// ParseSwiftCodes wczytuje i przetwarza dane z pliku CSV
func ParseSwiftCodes(filePath string) ([]models.SwiftCode, error) {
	// 1. Otwieramy plik CSV
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 2. Tworzymy czytnik CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %v", err)
	}

	// 3. Przygotowujemy listę wynikową i mapę siedzib głównych
	var swiftCodes []models.SwiftCode
	headquarterMap := make(map[string]string)

	// 4. Przetwarzamy każdy rekord z CSV
	for i, record := range records {
		if i == 0 {
			// Pomijamy nagłówki
			continue
		}

		// 5. Parsujemy dane z rekordu
		swiftCode := strings.ToUpper(record[1]) // Kolumna B: SWIFT CODE
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

		// 6. Dodajemy do mapy, jeśli to siedziba główna
		if isHeadquarter {
			headquarterMap[swiftCode[:8]] = swiftCode
		}

		// 7. Tworzymy instancję SwiftCode z przetworzonymi danymi
		swift := models.SwiftCode{
			SwiftCode:     swiftCode,
			BankName:      strings.TrimSpace(record[3]), // Kolumna D: NAME
			Address:       strings.TrimSpace(record[4]), // Kolumna E: ADDRESS
			TownName:      strings.TrimSpace(record[5]), // Kolumna F: TOWN NAME
			CountryISO2:   strings.ToUpper(record[0]),   // Kolumna A: COUNTRY ISO2 CODE
			CountryName:   strings.TrimSpace(record[6]), // Kolumna G: COUNTRY NAME
			TimeZone:      strings.TrimSpace(record[7]), // Kolumna H: TIME ZONE
			IsHeadquarter: isHeadquarter,
		}

		// 8. Powiązanie oddziału z siedzibą główną
		if !isHeadquarter {
			if hq, ok := headquarterMap[swiftCode[:8]]; ok {
				swift.HeadquarterID = &hq
			}
		}

		// 9. Dodajemy rekord do listy wynikowej
		swiftCodes = append(swiftCodes, swift)
	}
	return swiftCodes, nil
}
