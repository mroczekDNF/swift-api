package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/tests/mocks" // Importuj mocki z katalogu mocks
	"github.com/stretchr/testify/assert"
)

func TestGetSwiftCodeDetails_Integration(t *testing.T) {
	// Przygotowanie routera i mock repozytorium
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockSwiftCodeRepository) // Użycie mocka z katalogu mocks
	handler := handlers.NewSwiftCodeHandler(mockRepo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	// Przykładowy SWIFT code i dane z mocka
	swiftCode := "BANKUS33XXX"
	mockSwift := &models.SwiftCode{
		SwiftCode:     swiftCode,
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		Address:       "123 Test St",
	}

	// Konfiguracja mocka
	mockRepo.On("GetBySwiftCode", swiftCode).Return(mockSwift, nil)
	mockRepo.On("GetBranchesByHeadquarter", swiftCode).Return([]models.SwiftCode{
		{
			SwiftCode:     "BANKUS33ABC",
			BankName:      "Test Bank Branch",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: false,
			Address:       "456 Test Ave",
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

	// Weryfikacja pól odpowiedzi
	assert.Equal(t, swiftCode, response["swiftCode"])
	assert.Equal(t, "Test Bank", response["bankName"])
	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])
	assert.True(t, response["isHeadquarter"].(bool))

	// Weryfikacja branchy
	branches := response["branches"].([]interface{})
	assert.Len(t, branches, 1)
	branch := branches[0].(map[string]interface{})
	assert.Equal(t, "BANKUS33ABC", branch["swiftCode"])
	assert.Equal(t, "Test Bank Branch", branch["bankName"])

	// Sprawdzamy, czy mock został użyty zgodnie z oczekiwaniami
	mockRepo.AssertExpectations(t)
}
