package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddSwiftCode_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)

	requestBody := map[string]interface{}{
		"swiftCode":     "NEWBANKXYY",
		"bankName":      "New Bank",
		"address":       "123 New St",
		"countryISO2":   "US",
		"countryName":   "United States",
		"isHeadquarter": true,
	}

	expectedSwiftCode := &models.SwiftCode{
		SwiftCode:     "NEWBANKXYY",
		BankName:      "New Bank",
		Address:       "123 New St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}

	mockRepo.On("GetBySwiftCode", "NEWBANKXYY").Return(nil, nil)
	mockRepo.On("InsertSwiftCode", expectedSwiftCode).Return(nil)

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "SWIFT code added successfully", response["message"])

	mockRepo.AssertExpectations(t)
}

func TestAddSwiftCode_InvalidRequestStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)

	invalidRequestBody := map[string]interface{}{
		"swiftCode": "BANK123XXX",
	}

	body, _ := json.Marshal(invalidRequestBody)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["error"], "Invalid request structure")
	mockRepo.AssertExpectations(t)
}

func TestAddSwiftCode_AlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)

	requestBody := map[string]interface{}{
		"swiftCode":     "EXISTBANKXX",
		"bankName":      "Existing Bank",
		"address":       "Old Address",
		"countryISO2":   "GB",
		"countryName":   "United Kingdom",
		"isHeadquarter": false,
	}

	existingSwiftCode := &models.SwiftCode{
		SwiftCode:     "EXISTBANKXX",
		BankName:      "Existing Bank",
		Address:       "Old Address",
		CountryISO2:   "GB",
		CountryName:   "United Kingdom",
		IsHeadquarter: false,
	}

	mockRepo.On("GetBySwiftCode", "EXISTBANKXX").Return(existingSwiftCode, nil)

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusConflict, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "SWIFT code already exists in the database", response["error"])

	mockRepo.AssertExpectations(t)
}
