package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		"swiftCode":     swift.SwiftCode,
		"bankName":      swift.BankName,
		"address":       swift.Address,
		"countryISO2":   swift.CountryISO2,
		"countryName":   swift.CountryName,
		"isHeadquarter": swift.IsHeadquarter,
	}

	// Jeśli to headquarters, dodaj branches do odpowiedzi
	if swift.IsHeadquarter {
		branches, err := h.repo.GetBranchesByHeadquarter(swift.SwiftCode)
		if err != nil && err != sql.ErrNoRows {
			log.Println("Błąd pobierania branchy:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania branchy"})
			return
		}

		// **Fix: Konwersja branchy do poprawnej struktury JSON**
		branchList := make([]gin.H, 0)
		for _, branch := range branches {
			branchList = append(branchList, gin.H{
				"swiftCode":   branch.SwiftCode,
				"bankName":    branch.BankName,
				"address":     branch.Address,
				"countryISO2": branch.CountryISO2,
				"countryName": branch.CountryName,
			})
		}
		response["branches"] = branchList
	}

	c.JSON(http.StatusOK, response)
}
