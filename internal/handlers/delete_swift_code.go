package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/models"
)

// DeleteSwiftCode obsługuje żądanie DELETE /v1/swift-codes/{swift-code}
func DeleteSwiftCode(c *gin.Context) {
	// Pobierz swift-code z parametrów URL
	swiftCode := strings.ToUpper(strings.TrimSpace(c.Param("swift-code")))

	// Rozpocznij transakcję
	tx := db.DB.Begin()

	// Wyszukaj rekord w bazie danych
	var swift models.SwiftCode
	if err := tx.Where("swift_code = ?", swiftCode).First(&swift).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"message": "SWIFT code not found"})
		return
	}

	// Jeśli rekord jest headquarterem, usuń powiązania branchy
	if swift.IsHeadquarter {
		if err := tx.Model(&models.SwiftCode{}).
			Where("headquarter_id = ?", swift.SwiftCode).
			Update("headquarter_id", nil).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update branches", "details": err.Error()})
			return
		}
	}

	// Usuń rekord
	if err := tx.Delete(&swift).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete SWIFT code", "details": err.Error()})
		return
	}

	// Zatwierdź transakcję
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Transaction failed", "details": err.Error()})
		return
	}

	// Zwróć odpowiedź
	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully"})
}
