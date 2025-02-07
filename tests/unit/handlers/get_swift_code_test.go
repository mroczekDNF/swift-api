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

func TestGetSwiftCodeDetails_Headquarter(t *testing.T) {
	// Prepare the router and mock repository
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	// Example SWIFT code for headquarters
	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		Address:       "123 Test St",
		HeadquarterID: nil, // Headquarters should have nil
	}

	// Mock configuration for headquarters
	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("GetBranchesByHeadquarter", swiftCode).Return([]models.SwiftCode{
		{
			ID:            2,
			SwiftCode:     "BANKUS33ABC",
			BankName:      "Test Bank Branch",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: false,
			Address:       "456 Test Ave",
			HeadquarterID: &mockSwift.ID, // Points to headquarters
		},
	}, nil)

	// Test request
	req, _ := http.NewRequest("GET", "/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	// Check the response
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	fmt.Printf("DEBUG: API Response: %s\n", recorder.Body.String())

	// Verify the fields for headquarters
	assert.Equal(t, swiftCode, response["swiftCode"])
	assert.Equal(t, "Test Bank", response["bankName"])
	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])
	assert.True(t, response["isHeadquarter"].(bool))

	// Verify branches
	branches, branchesExist := response["branches"].([]interface{})
	assert.True(t, branchesExist, "branches field should exist for headquarters")
	assert.Len(t, branches, 1)

	// Fix branch assertions
	branchData, ok := branches[0].(map[string]interface{})
	assert.True(t, ok, "Branch data should be a map[string]interface{}")
	assert.Equal(t, "BANKUS33ABC", branchData["swiftCode"])
	assert.Equal(t, "Test Bank Branch", branchData["bankName"])
	assert.Equal(t, "US", branchData["countryISO2"])
	assert.Equal(t, "United States", branchData["countryName"])
	assert.Equal(t, "456 Test Ave", branchData["address"])

	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodeDetails_Headquarter_NoBranches(t *testing.T) {
	// Check if a headquarters can exist without branches (branches == [])

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		Address:       "123 Test St",
		HeadquarterID: nil,
	}

	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("GetBranchesByHeadquarter", swiftCode).Return([]models.SwiftCode{}, nil) // No branches

	req, _ := http.NewRequest("GET", "/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, swiftCode, response["swiftCode"])
	assert.Equal(t, "Test Bank", response["bankName"])
	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])
	assert.True(t, response["isHeadquarter"].(bool))

	// Should return an empty list of branches
	branches, branchesExist := response["branches"].([]interface{})
	assert.True(t, branchesExist, "branches field should exist for headquarters")
	assert.Len(t, branches, 0) // No branches

	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodeDetails_Branch(t *testing.T) {
	// Example SWIFT code for a branch
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	swiftCode := "BANKUS33ABC"
	headquarterID := int64(1)
	mockSwift := &models.SwiftCode{
		ID:            2,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank Branch",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: false,
		Address:       "456 Test Ave",
		HeadquarterID: &headquarterID,
	}

	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)

	req, _ := http.NewRequest("GET", "/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, swiftCode, response["swiftCode"])
	assert.Equal(t, "Test Bank Branch", response["bankName"])
	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])
	assert.False(t, response["isHeadquarter"].(bool))

	_, branchesExist := response["branches"]
	assert.False(t, branchesExist, "branches field should NOT exist for branch SWIFT code")

	mockRepo.AssertExpectations(t)
}
