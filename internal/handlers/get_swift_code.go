package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// GetSwiftCodeDetails zwraca szczegóły dla danego SWIFT code.
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

	// Struktura odpowiedzi
	response := gin.H{
		"address":       swift.Address,
		"bankName":      swift.BankName,
		"countryISO2":   swift.CountryISO2,
		"countryName":   swift.CountryName,
		"isHeadquarter": swift.IsHeadquarter,
		"swiftCode":     swift.SwiftCode,
	}

	// Jeśli to headquarters, dodaj branches do odpowiedzi
	if swift.IsHeadquarter {
		branches, err := h.repo.GetBranchesByHeadquarter(swift.SwiftCode)
		if err != nil && err != sql.ErrNoRows { // Jeśli to poważny błąd, zwracamy 500
			log.Println("Błąd pobierania branchy:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania branchy"})
			return
		}

		// Jeśli brak branchy, zwracamy pustą listę
		if len(branches) == 0 {
			branches = []models.SwiftCode{}
		}
		response["branches"] = branches
	}

	c.JSON(http.StatusOK, response)
}
