package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetSwiftCodeDetails returns details for a given SWIFT code.
func (h *SwiftCodeHandler) GetSwiftCodeDetails(c *gin.Context) {
	swiftCode := strings.TrimSpace(c.Param("swiftCode"))

	// Fetch details for the given SWIFT code
	swift, err := h.repo.GetBySwiftCode(swiftCode)
	if err != nil {
		log.Println("Error fetching SWIFT code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}
	if swift == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
		return
	}

	// Response structure
	response := gin.H{
		"swiftCode":     swift.SwiftCode,
		"bankName":      swift.BankName,
		"address":       swift.Address,
		"countryISO2":   swift.CountryISO2,
		"countryName":   swift.CountryName,
		"isHeadquarter": swift.IsHeadquarter,
	}

	// If it's a headquarters, add branches to the response
	if swift.IsHeadquarter {
		branches, err := h.repo.GetBranchesByHeadquarter(swift.SwiftCode)
		if err != nil && err != sql.ErrNoRows {
			log.Println("Error fetching branches:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching branches"})
			return
		}

		// Convert branches to the proper JSON structure
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
