package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/handlers"
	"github.com/mroczekDNF/swift-api/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestAddSwiftCode_Success_HQ(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)

	body := gin.H{
		"swiftCode":     "BANKFR77XXX",
		"bankName":      "Credit Agricole HQ",
		"address":       "456 Paris Ave",
		"countryISO2":   "FR",
		"countryName":   "France",
		"isHeadquarter": true,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SWIFT code added successfully", response["message"])

	swiftCode, _ := repo.GetBySwiftCode("BANKFR55XXX")
	assert.NotNil(t, swiftCode)
	assert.Equal(t, "BANKFR55XXX", swiftCode.SwiftCode)
	assert.Equal(t, "Credit Agricole HQ", swiftCode.BankName)
	assert.Equal(t, "456 Paris Ave", swiftCode.Address)
	assert.Equal(t, "FR", swiftCode.CountryISO2)
	assert.Equal(t, "France", swiftCode.CountryName)
	assert.True(t, swiftCode.IsHeadquarter)
}

func TestAddSwiftCode_Success_Branch_WithHQ(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)

	body := gin.H{
		"swiftCode":     "BANKCA66AAA",
		"bankName":      "Canadian Bank Branch A",
		"address":       "456 Nice Blvd",
		"countryISO2":   "CA",
		"countryName":   "Canada",
		"isHeadquarter": false,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SWIFT code added successfully", response["message"])

	swiftCode, err := repo.GetBySwiftCode("BANKCA66AAA")
	assert.NoError(t, err)
	assert.NotNil(t, swiftCode)
	assert.Equal(t, "BANKCA66AAA", swiftCode.SwiftCode)
	assert.Equal(t, "Canadian Bank Branch A", swiftCode.BankName)
	assert.Equal(t, "456 Nice Blvd", swiftCode.Address)
	assert.Equal(t, "CA", swiftCode.CountryISO2)
	assert.Equal(t, "Canada", swiftCode.CountryName)
	assert.False(t, swiftCode.IsHeadquarter)

	assert.NotNil(t, swiftCode.HeadquarterID)

	hq, err := repo.GetBySwiftCode("BANKCA66XXX")
	assert.NoError(t, err)
	assert.NotNil(t, hq)
	assert.Equal(t, hq.ID, *swiftCode.HeadquarterID)
}

func TestAddSwiftCode_Conflict(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)

	body := gin.H{
		"swiftCode":     "BANKUS33XXX",
		"bankName":      "Test Bank USA HQ",
		"address":       "123 Test St",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusConflict, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "SWIFT code already exists in the database", response["error"])
}
