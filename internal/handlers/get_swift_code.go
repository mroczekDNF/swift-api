package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// GetSwiftCodeDetails zwraca szczegóły dla danego SWIFT code
func GetSwiftCodeDetails(c *gin.Context) {
	swiftCode := strings.TrimSpace(c.Param("swift-code")) // Pobierz {swift-code} z URL

	log.Printf("Param swift-code: %s", c.Param("swift-code"))

	// Pobierz szczegóły dla podanego SWIFT code
	var swift models.SwiftCode
	if err := db.DB.Where("swift_code = ?", swiftCode).First(&swift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
		return
	}

	// Jeśli to headquarter, znajdź wszystkie branche
	if swift.IsHeadquarter {
		var branches []models.SwiftCode
		if err := db.DB.Where("headquarter_id = ?", swift.SwiftCode).Find(&branches).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branches"})
			return
		}

		// Struktura odpowiedzi dla headquarter
		response := gin.H{
			"address":       swift.Address,
			"bankName":      swift.BankName,
			"countryISO2":   swift.CountryISO2,
			"countryName":   swift.CountryName,
			"isHeadquarter": swift.IsHeadquarter,
			"swiftCode":     swift.SwiftCode,
			"branches":      branches, // Lista branchy
		}

		c.JSON(http.StatusOK, response)
		return
	}

	// Jeśli to branch, zwróć tylko jego szczegóły
	response := gin.H{
		"address":       swift.Address,
		"bankName":      swift.BankName,
		"countryISO2":   swift.CountryISO2,
		"countryName":   swift.CountryName,
		"isHeadquarter": swift.IsHeadquarter,
		"swiftCode":     swift.SwiftCode,
	}

	c.JSON(http.StatusOK, response)
}
