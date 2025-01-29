package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// GetSwiftCodeDetails zwraca szczegóły dla danego SWIFT code
func (h *SwiftCodeHandler) GetSwiftCodeDetails(c *gin.Context) {
	swiftCode := strings.TrimSpace(c.Param("swiftCode"))

	// Pobierz szczegóły dla podanego SWIFT code
	swift, err := h.repo.GetBySwiftCode(swiftCode)
	if err != nil {
		log.Println("Błąd pobierania SWIFT code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania danych"})
		return
	}
	if swift == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code nie znaleziony"})
		return
	}

	// Jeśli to headquarter, znajdź wszystkie branche
	var branches []models.SwiftCode
	if swift.IsHeadquarter {
		branches, err = h.repo.GetBranchesByHeadquarter(swift.SwiftCode)
		if err != nil {
			log.Println("Błąd pobierania branchy:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania branchy"})
			return
		}
	}

	// Struktura odpowiedzi
	response := gin.H{
		"address":       swift.Address,
		"bankName":      swift.BankName,
		"countryISO2":   swift.CountryISO2,
		"countryName":   swift.CountryName,
		"isHeadquarter": swift.IsHeadquarter,
		"swiftCode":     swift.SwiftCode,
		"branches":      branches,
	}

	c.JSON(http.StatusOK, response)
}
