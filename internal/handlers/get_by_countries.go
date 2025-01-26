package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// GetSwiftCodesByCountry zwraca wszystkie SWIFT codes dla danego kraju
func GetSwiftCodesByCountry(c *gin.Context) {
	// Pobierz kod ISO-2 kraju z parametrów URL
	countryISO2 := strings.TrimSpace(c.Param("countryISO2code"))

	// Znajdź wszystkie rekordy dla danego kraju
	var swiftCodes []models.SwiftCode
	if err := db.DB.Where("country_iso2 = ?", countryISO2).Find(&swiftCodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SWIFT codes"})
		return
	}

	// Jeśli nie znaleziono żadnych rekordów
	if len(swiftCodes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No SWIFT codes found for the specified country"})
		return
	}

	// Ustaw kraj na podstawie pierwszego rekordu
	countryName := swiftCodes[0].CountryName

	// Tworzenie odpowiedzi
	response := gin.H{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  swiftCodes, // Lista SWIFT codes
	}

	c.JSON(http.StatusOK, response)
}
