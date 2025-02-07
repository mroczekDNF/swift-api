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

func TestGetSwiftCodeDetails_Headquarter(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	req, _ := http.NewRequest("GET", "/swift-codes/BANKUS33XXX", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	log.Printf("Checking JSON keys: %+v", response)

	assert.Equal(t, "BANKUS33XXX", response["swiftCode"])
	assert.Equal(t, "Test Bank USA", response["bankName"])
	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "United States", response["countryName"])
	assert.True(t, response["isHeadquarter"].(bool))

	branches, branchesExist := response["branches"].([]interface{})
	assert.True(t, branchesExist, "The 'branches' field should exist")
	assert.Len(t, branches, 2, "There should be exactly 2 branches")

	log.Println("Branch data in response:")
	expectedBranches := map[string]string{
		"BANKUS33ABC": "Test Bank Branch A",
		"BANKUS33DEF": "Test Bank Branch B",
	}

	for _, branch := range branches {
		branchData := branch.(map[string]interface{})
		swiftCode, _ := branchData["swiftCode"].(string)
		bankName, _ := branchData["bankName"].(string)

		expectedName, exists := expectedBranches[swiftCode]
		assert.True(t, exists, "Unexpected branch: "+swiftCode)
		assert.Equal(t, expectedName, bankName)
	}
}

func TestGetSwiftCodeDetails_HQWithoutBranches(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	req, _ := http.NewRequest("GET", "/swift-codes/BANKDE44XXX", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "BANKDE44XXX", response["swiftCode"])
	assert.Equal(t, "Deutsche Bank HQ", response["bankName"])
	assert.Equal(t, "DE", response["countryISO2"])
	assert.Equal(t, "Germany", response["countryName"])
	assert.True(t, response["isHeadquarter"].(bool))

	// Verify that 'branches' field exists and is an empty list
	branches, branchesExist := response["branches"]
	assert.True(t, branchesExist, "The 'branches' field should exist")
	assert.Len(t, branches, 0, "The 'branches' field should be an empty list")
}

func TestGetSwiftCodeDetails_BranchOnly(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.GET("/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)

	req, _ := http.NewRequest("GET", "/swift-codes/BANKGB22AAA", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	log.Printf("JSON Response: %s", recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "BANKGB22AAA", response["swiftCode"])
	assert.Equal(t, "Independent Branch UK", response["bankName"])
	assert.Equal(t, "GB", response["countryISO2"])
	assert.Equal(t, "United Kingdom", response["countryName"])
	assert.False(t, response["isHeadquarter"].(bool))
	_, branchesExist := response["branches"].([]interface{})
	assert.False(t, branchesExist, "The 'branches' field should not exist for a branch")
}
