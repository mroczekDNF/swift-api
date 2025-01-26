package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// AddSwiftCode obsługuje żądanie POST /v1/swift-codes/
func AddSwiftCode(c *gin.Context) {
	// Struktura dla żądania
	var request struct {
		Address       string `json:"address" binding:"required"`
		BankName      string `json:"bankName" binding:"required"`
		CountryISO2   string `json:"countryISO2" binding:"required"`
		CountryName   string `json:"countryName" binding:"required"`
		IsHeadquarter *bool  `json:"isHeadquarter" binding:"required"`
		SwiftCode     string `json:"swiftCode" binding:"required"`
	}

	// Walidacja JSON-a
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request structure", "details": err.Error()})
		return
	}

	// Upewnij się, że `isHeadquarter` nie jest `nil`
	if request.IsHeadquarter == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'isHeadquarter' is required"})
		return
	}

	// Logowanie wartości IsHeadquarter
	log.Printf("IsHeadquarter value: %v", *request.IsHeadquarter)

	// Normalizacja danych
	swiftCode := strings.ToUpper(strings.TrimSpace(request.SwiftCode))
	countryISO2 := strings.ToUpper(strings.TrimSpace(request.CountryISO2))
	countryName := strings.ToUpper(strings.TrimSpace(request.CountryName))

	// Sprawdź, czy rekord z tym SWIFT code już istnieje
	var existing models.SwiftCode
	if err := db.DB.Where("swift_code = ?", swiftCode).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "SWIFT code already exists"})
		return
	}

	// Utwórz nowy rekord
	newSwiftCode := models.SwiftCode{
		SwiftCode:     swiftCode,
		BankName:      strings.TrimSpace(request.BankName),
		Address:       strings.TrimSpace(request.Address),
		CountryISO2:   countryISO2,
		CountryName:   countryName,
		IsHeadquarter: *request.IsHeadquarter,
	}

	// Jeśli to branch, znajdź odpowiedni headquarter
	// if !*request.IsHeadquarter {
	// 	var headquarter models.SwiftCode
	// 	if err := db.DB.Where("swift_code = ?", swiftCode[:8]+"XXX").First(&headquarter).Error; err == nil {
	// 		newSwiftCode.HeadquarterID = &headquarter.SwiftCode
	// 	} else {
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Headquarter not found for branch"})
	// 		return
	// 	}
	// }
	// tu sie trzeba zastanowic jeszcze
	if !*request.IsHeadquarter {
		var headquarter models.SwiftCode
		if err := db.DB.Where("swift_code = ?", swiftCode[:8]+"XXX").First(&headquarter).Error; err == nil {
			newSwiftCode.HeadquarterID = &headquarter.SwiftCode
		} else {
			// Logowanie informacyjne, jeśli nie znaleziono headquarter
			log.Printf("Headquarter not found for branch with SWIFT code: %s", swiftCode)
			// Pole HeadquarterID pozostanie puste (nil)
		}
	}
	// Zapisz rekord do bazy danych
	if err := db.DB.Create(&newSwiftCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save SWIFT code", "details": err.Error()})
		return
	}

	// Zwróć odpowiedź
	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code added successfully"})
}
