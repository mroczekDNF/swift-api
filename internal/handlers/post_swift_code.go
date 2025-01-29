package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// AddSwiftCode obsługuje żądanie POST /v1/swift-codes/
func (h *SwiftCodeHandler) AddSwiftCode(c *gin.Context) {
	// Struktura dla żądania
	var request struct {
		Address       string `json:"address" binding:"required"`
		BankName      string `json:"bankName" binding:"required"`
		CountryISO2   string `json:"countryISO2" binding:"required"`
		CountryName   string `json:"countryName" binding:"required"`
		IsHeadquarter *bool  `json:"isHeadquarter" binding:"required"`
		SwiftCode     string `json:"swiftCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niepoprawna struktura żądania", "details": err.Error()})
		return
	}

	if request.IsHeadquarter == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pole 'isHeadquarter' jest wymagane"})
		return
	}

	// Normalizacja danych
	swiftCode := strings.ToUpper(strings.TrimSpace(request.SwiftCode))
	countryISO2 := strings.ToUpper(strings.TrimSpace(request.CountryISO2))
	countryName := strings.ToUpper(strings.TrimSpace(request.CountryName))

	// Sprawdź, czy rekord z tym SWIFT code już istnieje
	exists, err := h.repo.GetBySwiftCode(swiftCode)
	if err != nil {
		log.Println("Błąd sprawdzania istnienia SWIFT code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Błąd sprawdzania danych"})
		return
	}
	if exists != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "SWIFT code już istnieje"})
		return
	}

	// Tworzenie nowego rekordu
	newSwiftCode := models.SwiftCode{
		SwiftCode:     swiftCode,
		BankName:      strings.TrimSpace(request.BankName),
		Address:       strings.TrimSpace(request.Address),
		CountryISO2:   countryISO2,
		CountryName:   countryName,
		IsHeadquarter: *request.IsHeadquarter,
	}

	// Jeśli to branch, znajdź odpowiedni headquarter
	if !*request.IsHeadquarter {
		headquarter, err := h.repo.GetBySwiftCode(swiftCode[:8] + "XXX")
		if err != nil {
			log.Println("Błąd wyszukiwania headquarter:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Błąd wyszukiwania headquarter"})
			return
		}
		if headquarter != nil {
			newSwiftCode.HeadquarterID = &headquarter.ID
		}
	}

	// Zapisz rekord do bazy danych
	if err := h.repo.InsertSwiftCode(&newSwiftCode); err != nil {
		log.Println("Błąd zapisu SWIFT code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Błąd zapisu SWIFT code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code dodany poprawnie"})
}
