package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// GetSwiftCodesByCountry returns SWIFT codes in a new response format
func (h *SwiftCodeHandler) GetSwiftCodesByCountry(c *gin.Context) {
	countryISO2 := strings.ToUpper(strings.TrimSpace(c.Param("countryISO2")))

	swiftCodes, err := h.repo.GetByCountryISO2(countryISO2)
	if err != nil {
		log.Println("Error fetching SWIFT codes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching SWIFT codes"})
		return
	}

	if len(swiftCodes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No SWIFT codes found for the given country"})
		return
	}

	response := formatSwiftCodesResponse(countryISO2, swiftCodes)
	c.JSON(http.StatusOK, response)
}

// formatSwiftCodesResponse formats the SWIFT codes into the expected response structure
func formatSwiftCodesResponse(countryISO2 string, swiftCodes []models.SwiftCode) gin.H {
	countryName := swiftCodes[0].CountryName

	var formattedSwiftCodes []gin.H
	for _, code := range swiftCodes {
		formattedSwiftCodes = append(formattedSwiftCodes, gin.H{
			"address":       code.Address,
			"bankName":      code.BankName,
			"countryISO2":   code.CountryISO2,
			"isHeadquarter": code.IsHeadquarter,
			"swiftCode":     code.SwiftCode,
		})
	}

	return gin.H{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  formattedSwiftCodes,
	}
}
