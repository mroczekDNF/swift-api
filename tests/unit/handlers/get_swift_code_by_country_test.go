package unit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetSwiftCodesByCountry_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	countryISO2 := "US"
	mockSwiftCodes := []models.SwiftCode{
		{
			ID:            1,
			SwiftCode:     "BANKUS33XXX",
			BankName:      "Test Bank HQ",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: true,
			Address:       "123 Test HQ Street",
		},
		{
			ID:            2,
			SwiftCode:     "BANKUS33ABC",
			BankName:      "Test Bank Branch",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: false,
			Address:       "456 Test Branch Avenue",
		},
	}

	mockRepo.On("GetByCountryISO2", countryISO2).Return(mockSwiftCodes, nil)

	req, _ := http.NewRequest("GET", "/swift-codes/country/"+countryISO2, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, countryISO2, response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])

	swiftCodes, ok := response["swiftCodes"].([]interface{})
	assert.True(t, ok, "swiftCodes field should exist and be a list")
	assert.Len(t, swiftCodes, 2)

	// Validate first SWIFT code
	code1 := swiftCodes[0].(map[string]interface{})
	assert.Equal(t, "BANKUS33XXX", code1["swiftCode"])
	assert.Equal(t, "Test Bank HQ", code1["bankName"])
	assert.Equal(t, "123 Test HQ Street", code1["address"])
	assert.True(t, code1["isHeadquarter"].(bool))

	// Validate second SWIFT code
	code2 := swiftCodes[1].(map[string]interface{})
	assert.Equal(t, "BANKUS33ABC", code2["swiftCode"])
	assert.Equal(t, "Test Bank Branch", code2["bankName"])
	assert.Equal(t, "456 Test Branch Avenue", code2["address"])
	assert.False(t, code2["isHeadquarter"].(bool))

	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodesByCountry_NoRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	countryISO2 := "XX"
	mockRepo.On("GetByCountryISO2", countryISO2).Return([]models.SwiftCode{}, nil)

	req, _ := http.NewRequest("GET", "/swift-codes/country/"+countryISO2, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["error"], "No SWIFT codes found for the given country")
	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodesByCountry_NoSwiftCodesFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	countryISO2 := "XX" // Country ISO2 code that does not exist in the mock database

	// Mock behavior: return an empty list when queried with an invalid or missing country
	mockRepo.On("GetByCountryISO2", countryISO2).Return([]models.SwiftCode{}, nil)

	req, _ := http.NewRequest("GET", "/swift-codes/country/"+countryISO2, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Expect a 404 Not Found response
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify the error message in the response
	assert.Contains(t, response["error"], "No SWIFT codes found for the given country")

	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodesByCountry_RepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)

	countryISO2 := "US"
	mockRepo.On("GetByCountryISO2", countryISO2).Return(nil, fmt.Errorf("database error"))

	req, _ := http.NewRequest("GET", "/swift-codes/country/"+countryISO2, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["error"], "Error fetching SWIFT codes")
	mockRepo.AssertExpectations(t)
}
