package unit

import (
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

func TestDeleteSwiftCode_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank HQ",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: false,
		Address:       "123 Test HQ Street",
	}

	// Mock behavior
	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("DeleteSwiftCode", swiftCode).Return(nil)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "SWIFT code deleted successfully")
	mockRepo.AssertExpectations(t)
}

func TestDeleteSwiftCode_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "INVALIDCODE"

	// Mock behavior for non-existing SWIFT code
	mockRepo.On("GetBySwiftCode", swiftCode).Return(nil, nil)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "SWIFT code not found")
	mockRepo.AssertExpectations(t)
}

func TestDeleteSwiftCode_ErrorRetrievingCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKUS33XXX"

	// Simulate a database error when retrieving the SWIFT code
	mockRepo.On("GetBySwiftCode", swiftCode).Return(nil, assert.AnError)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "Error retrieving SWIFT code")
	mockRepo.AssertExpectations(t)
}

func TestDeleteSwiftCode_ErrorDeleting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank HQ",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: false,
		Address:       "123 Test HQ Street",
	}

	// Simulate a successful retrieval but failure on delete
	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("DeleteSwiftCode", swiftCode).Return(assert.AnError)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "Error deleting SWIFT code")
	mockRepo.AssertExpectations(t)
}

func TestDeleteSwiftCode_HeadquarterWithBranches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank HQ",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		Address:       "123 Test HQ Street",
	}

	// Simulate successful retrieval and detachment of branches
	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("DetachBranchesFromHeadquarter", mockSwift.ID).Return(nil)
	mockRepo.On("DeleteSwiftCode", swiftCode).Return(nil)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "SWIFT code deleted successfully")
	mockRepo.AssertExpectations(t)
}

func TestDeleteSwiftCode_ErrorDetachingBranches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank HQ",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		Address:       "123 Test HQ Street",
	}

	// Simulate a successful retrieval but failure when detaching branches
	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("DetachBranchesFromHeadquarter", mockSwift.ID).Return(assert.AnError)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "Error detaching branches")
	mockRepo.AssertExpectations(t)
}
