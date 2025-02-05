package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// Globalne wyrażenia regularne dla optymalizacji
var (
	// Kod kraju: dokładnie 2 wielkie litery.
	countryISO2Regex = regexp.MustCompile(`^[A-Z]{2}$`)
	// SWIFT: 4 litery (banku), 2 litery (kraju), 2 alfanumeryczne (lokalizacji), opcjonalnie 3 alfanumeryczne (oddziału).
	swiftCodeRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
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

	// Sprawdzenie czy pole IsHeadquarter nie jest nil.
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
	if err := assignHeadquarterID(h, &newSwiftCode); err != nil {
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
	// Sprawdzamy długość kodu SWIFT.
	if len(request.SwiftCode) < 8 || len(request.SwiftCode) > 11 {
		return &ValidationError{"Invalid SWIFT code length. Must be between 8 and 11 characters."}
	}

	// Walidacja formatu SWIFT
	if !swiftCodeRegex.MatchString(request.SwiftCode) {
		return &ValidationError{"Invalid SWIFT code format. Expected format: 4 letters, 2 letters, 2 alphanumeric characters and optional 3 alphanumeric characters."}
	}

	// Walidacja kodu kraju
	if len(request.CountryISO2) != 2 || !countryISO2Regex.MatchString(request.CountryISO2) {
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

// assignHeadquarterID ustawia HeadquarterID dla SwiftCode, jeśli dany kod nie jest headquarter.
// Wyszukuje headquarter przy pomocy pierwszych 8 znaków SwiftCode, uzupełniając końcówkę "XXX".
func assignHeadquarterID(h *SwiftCodeHandler, newSwiftCode *models.SwiftCode) error {
	if !newSwiftCode.IsHeadquarter {
		headquarterCode := newSwiftCode.SwiftCode[:8] + "XXX"
		headquarter, err := h.repo.GetBySwiftCode(headquarterCode)
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
