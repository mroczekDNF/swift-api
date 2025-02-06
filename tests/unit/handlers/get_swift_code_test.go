package integration

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
	// Przygotowanie routera i mock repozytorium
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository)
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	// Przykładowy SWIFT code dla headquarters
	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		ID:            1,
		SwiftCode:     swiftCode,
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		Address:       "123 Test St",
		HeadquarterID: nil, // Headquarters powinien mieć nil
	}

	// Konfiguracja mocka dla headquarters
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
			HeadquarterID: &mockSwift.ID, // Wskazuje na headquarters
		},
	}, nil)

	// Testowe żądanie
	req, _ := http.NewRequest("GET", "/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	// Sprawdzanie odpowiedzi
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	fmt.Printf("DEBUG: Odpowiedź API: %s\n", recorder.Body.String())

	// Weryfikacja pól odpowiedzi dla headquarters
	assert.Equal(t, swiftCode, response["swiftCode"])
	assert.Equal(t, "Test Bank", response["bankName"])
	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])
	assert.True(t, response["isHeadquarter"].(bool))

	// // Weryfikacja branchy
	branches, branchesExist := response["branches"].([]interface{})
	assert.True(t, branchesExist, "branches field should exist for headquarters")
	assert.Len(t, branches, 1)

	branch := branches[0].(map[string]interface{})
	assert.Equal(t, "BANKUS33ABC", branch["SwiftCode"])
	assert.Equal(t, "Test Bank Branch", branch["BankName"])

	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodeDetails_Headquarter_NoBranches(t *testing.T) {
	// Sprawdzamy, czy headquarters może istnieć bez branchy (branches == [])

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
	mockRepo.On("GetBranchesByHeadquarter", swiftCode).Return([]models.SwiftCode{}, nil) // Brak branchy

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

	// Powinno zwracać pustą listę branchy
	branches, branchesExist := response["branches"].([]interface{})
	assert.True(t, branchesExist, "branches field should exist for headquarters")
	assert.Len(t, branches, 0) // Brak branchy

	mockRepo.AssertExpectations(t)
}

func TestGetSwiftCodeDetails_Branch(t *testing.T) {
	// Przykładowy SWIFT code dla branch'a
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
		HeadquarterID: &headquarterID, // Branch wskazuje na headquarters
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

	// Branch NIE powinien zawierać "branches"
	_, branchesExist := response["branches"]
	assert.False(t, branchesExist, "branches field should NOT exist for branch SWIFT code")

	mockRepo.AssertExpectations(t)
}
