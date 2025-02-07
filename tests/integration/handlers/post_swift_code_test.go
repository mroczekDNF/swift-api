package integration

// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/mroczekDNF/swift-api/internal/db"
// 	"github.com/mroczekDNF/swift-api/internal/handlers"
// 	"github.com/mroczekDNF/swift-api/internal/repositories"
// 	"github.com/stretchr/testify/assert"
// )

// func TestAddSwiftCode_DefaultAddress(t *testing.T) {
// 	SetupTestDatabase(t)
// 	defer CleanupTestDatabase(t)

// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()

// 	repo := repositories.NewSwiftCodeRepository(db.DB)
// 	handler := handlers.NewSwiftCodeHandler(repo)
// 	router.POST("/swift-codes", handler.AddSwiftCode)

// 	requestBody := map[string]interface{}{
// 		"swiftCode":     "DEFAULT01XXX",
// 		"bankName":      "Bank Without Address",
// 		"address":       "",
// 		"countryISO2":   "FR",
// 		"countryName":   "France",
// 		"isHeadquarter": true,
// 	}

// 	body, _ := json.Marshal(requestBody)
// 	req, _ := http.NewRequest("POST", "/swift-codes", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	recorder := httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)

// 	log.Printf("JSON Response: %s", recorder.Body.String())

// 	assert.Equal(t, http.StatusOK, recorder.Code)

// 	var response map[string]interface{}
// 	err := json.Unmarshal(recorder.Body.Bytes(), &response)
// 	assert.NoError(t, err)

// 	assert.Equal(t, "SWIFT code added successfully", response["message"])

// 	// Sprawdź czy adres jest ustawiony jako "UNKNOWN"
// 	swiftCodeRepo, err := repo.GetBySwiftCode("DEFAULT01XXX")
// 	assert.NoError(t, err)                            // Upewnij się, że nie było błędu podczas pobierania
// 	assert.NotNil(t, swiftCodeRepo)                   // Upewnij się, że rekord istnieje
// 	assert.Equal(t, "UNKNOWN", swiftCodeRepo.Address) // Sprawdź, czy adres to "UNKNOWN"
// }
