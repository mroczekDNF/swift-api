package integration

import (
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

func TestGetSwiftCodesByCountry_Found(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	req, _ := http.NewRequest("GET", "/swift-codes/country/US", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])

	swiftCodes, swiftCodesExist := response["swiftCodes"].([]interface{})
	assert.True(t, swiftCodesExist, "The 'swiftCodes' field should exist")
	assert.GreaterOrEqual(t, len(swiftCodes), 1, "There should be at least one SWIFT code")

	firstCode := swiftCodes[0].(map[string]interface{})
	assert.NotEmpty(t, firstCode["swiftCode"], "Swift code should not be empty")
	assert.NotEmpty(t, firstCode["bankName"], "Bank name should not be empty")
	assert.Equal(t, "US", firstCode["countryISO2"])
}

func TestGetSwiftCodesByCountry_NotFound(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	req, _ := http.NewRequest("GET", "/swift-codes/country/ZZ", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected HTTP 500 due to database error handling")

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected JSON response to be parsable")

	assert.Contains(t, response, "error", "Expected 'error' field in the response")
	assert.Equal(t, "Error fetching SWIFT codes", response["error"], "Expected database error message")
}

func TestGetSwiftCodesByCountry_InvalidISO2(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	req, _ := http.NewRequest("GET", "/swift-codes/country/USA", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected HTTP 500 due to database error handling")

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected JSON response to be parsable")

	assert.Contains(t, response, "error", "Expected 'error' field in the response")
	assert.Equal(t, "Error fetching SWIFT codes", response["error"], "Expected database error message")
}
