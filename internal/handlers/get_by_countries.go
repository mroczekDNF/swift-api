package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetSwiftCodesByCountry zwraca wszystkie SWIFT codes dla danego kraju
func (h *SwiftCodeHandler) GetSwiftCodesByCountry(c *gin.Context) {
	countryISO2 := strings.ToUpper(strings.TrimSpace(c.Param("countryISO2")))

	// Pobierz SWIFT codes dla danego kraju
	swiftCodes, err := h.repo.GetByCountryISO2(countryISO2)
	if err != nil {
		log.Println("Błąd pobierania SWIFT codes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania SWIFT codes"})
		return
	}

	// Jeśli nie znaleziono żadnych rekordów
	if len(swiftCodes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brak SWIFT codes dla podanego kraju"})
		return
	}

	// Ustaw nazwę kraju na podstawie pierwszego rekordu
	countryName := swiftCodes[0].CountryName

	// Tworzenie odpowiedzi
	response := gin.H{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  swiftCodes,
	}

	c.JSON(http.StatusOK, response)
}
