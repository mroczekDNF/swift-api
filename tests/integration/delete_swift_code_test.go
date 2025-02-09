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

// Test deletion of an existing SWIFT Code (headquarter)
func TestDeleteSwiftCode_Success_HQ(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKDE44XXX" // Deutsche Bank HQ (Germany)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code deleted successfully", response["message"])
}

// Test deletion of a headquarter with branches
func TestDeleteSwiftCode_HQWithBranches(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	var branch1HQID, branch2HQID *int

	err := db.DB.QueryRow("SELECT headquarter_id FROM swift_codes WHERE swift_code = $1", "BANKUS33ABC").Scan(&branch1HQID)
	assert.NoError(t, err, "Failed to fetch headquarter_id for BANKUS33ABC")
	assert.NotNil(t, branch1HQID, "Expected BANKUS33ABC to have a headquarter_id before deletion")

	err = db.DB.QueryRow("SELECT headquarter_id FROM swift_codes WHERE swift_code = $1", "BANKUS33DEF").Scan(&branch2HQID)
	assert.NoError(t, err, "Failed to fetch headquarter_id for BANKUS33DEF")
	assert.NotNil(t, branch2HQID, "Expected BANKUS33DEF to have a headquarter_id before deletion")

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/BANKUS33XXX", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	log.Printf("JSON Response: %s", recorder.Body.String())

	var updatedBranch1HQID, updatedBranch2HQID *int

	err = db.DB.QueryRow("SELECT headquarter_id FROM swift_codes WHERE swift_code = $1", "BANKUS33ABC").Scan(&updatedBranch1HQID)
	assert.NoError(t, err, "Failed to fetch headquarter_id for BANKUS33ABC after deletion")
	assert.Nil(t, updatedBranch1HQID, "Expected BANKUS33ABC to have NULL headquarter_id after HQ deletion")

	err = db.DB.QueryRow("SELECT headquarter_id FROM swift_codes WHERE swift_code = $1", "BANKUS33DEF").Scan(&updatedBranch2HQID)
	assert.NoError(t, err, "Failed to fetch headquarter_id for BANKUS33DEF after deletion")
	assert.Nil(t, updatedBranch2HQID, "Expected BANKUS33DEF to have NULL headquarter_id after HQ deletion")
}

// Test deletion of an existing SWIFT Code (branch)
func TestDeleteSwiftCode_Success_Branch(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKCA77AAA" // Canadian Bank Branch A (Canada)

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code deleted successfully", response["message"])
}

// Test deletion of a non-existent SWIFT Code
func TestDeleteSwiftCode_NotFound(t *testing.T) {
	SetupTestDatabase(t)
	defer CleanupTestDatabase(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db.DB)
	handler := handlers.NewSwiftCodeHandler(repo)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	swiftCode := "BANKXX99ZZZ"

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+swiftCode, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var response map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code not found", response["message"])
}
