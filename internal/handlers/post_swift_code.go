package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
)

type SwiftCodeRequest struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName" binding:"required"`
	CountryISO2   string `json:"countryISO2" binding:"required"`
	CountryName   string `json:"countryName" binding:"required"`
	IsHeadquarter *bool  `json:"isHeadquarter" binding:"required"`
	SwiftCode     string `json:"swiftCode" binding:"required"`
}

func (h *SwiftCodeHandler) AddSwiftCode(c *gin.Context) {
	var request SwiftCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request structure", err.Error())
		return
	}

	if request.IsHeadquarter == nil {
		respondWithError(c, http.StatusBadRequest, "Field 'isHeadquarter' is required")
		return
	}

	normalizeSwiftCodeRequest(&request)

	if err := validateSwiftCodeRequest(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	existingCode, err := h.repo.GetBySwiftCode(request.SwiftCode)
	if err != nil {
		log.Println("Error checking if SWIFT code exists:", err)
		respondWithError(c, http.StatusInternalServerError, "Error checking data")
		return
	}
	if existingCode != nil {
		respondWithError(c, http.StatusConflict, "SWIFT code already exists in the database")
		return
	}

	newSwiftCode := createSwiftCodeModel(&request)
	if err := assignHeadquarterID(h, &newSwiftCode, request.SwiftCode); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Error finding headquarter")
		return
	}

	if err := h.repo.InsertSwiftCode(&newSwiftCode); err != nil {
		log.Println("Error saving SWIFT code:", err)
		respondWithError(c, http.StatusInternalServerError, "Error saving SWIFT code")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code added successfully"})
}

func normalizeSwiftCodeRequest(request *SwiftCodeRequest) {
	request.SwiftCode = strings.ToUpper(strings.TrimSpace(request.SwiftCode))
	request.CountryISO2 = strings.ToUpper(strings.TrimSpace(request.CountryISO2))
	request.CountryName = strings.TrimSpace(request.CountryName)
	request.BankName = strings.TrimSpace(request.BankName)
	request.Address = strings.TrimSpace(request.Address)

	if request.Address == "" {
		request.Address = "UNKNOWN"
	}
}

func validateSwiftCodeRequest(request *SwiftCodeRequest) error {
	if len(request.SwiftCode) < 8 || len(request.SwiftCode) > 11 {
		return &ValidationError{"Invalid SWIFT code length. Must be between 8 and 11 characters."}
	}

	if len(request.CountryISO2) != 2 || !regexp.MustCompile(`^[A-Z]{2}$`).MatchString(request.CountryISO2) {
		return &ValidationError{"Invalid country ISO2 code. Must be exactly 2 uppercase letters."}
	}

	if request.BankName == "" {
		return &ValidationError{"Bank name cannot be empty."}
	}

	if request.CountryName == "" {
		return &ValidationError{"Country name cannot be empty."}
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func createSwiftCodeModel(request *SwiftCodeRequest) models.SwiftCode {
	return models.SwiftCode{
		SwiftCode:     request.SwiftCode,
		BankName:      request.BankName,
		Address:       request.Address,
		CountryISO2:   request.CountryISO2,
		CountryName:   request.CountryName,
		IsHeadquarter: *request.IsHeadquarter,
	}
}

func assignHeadquarterID(h *SwiftCodeHandler, newSwiftCode *models.SwiftCode, swiftCode string) error {
	if !newSwiftCode.IsHeadquarter {
		headquarter, err := h.repo.GetBySwiftCode(swiftCode[:8] + "XXX")
		if err != nil {
			log.Println("Error finding headquarter:", err)
			return err
		}
		if headquarter != nil {
			newSwiftCode.HeadquarterID = &headquarter.ID
		}
	}
	return nil
}

func respondWithError(c *gin.Context, statusCode int, message string, details ...string) {
	errorResponse := gin.H{"error": message}
	if len(details) > 0 {
		errorResponse["details"] = details[0]
	}
	c.JSON(statusCode, errorResponse)
}
