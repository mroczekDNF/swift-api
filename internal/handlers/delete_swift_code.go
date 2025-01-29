package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/repositories"
)

// SwiftCodeHandler obsługuje operacje na kodach SWIFT
type SwiftCodeHandler struct {
	repo *repositories.SwiftCodeRepository
}

// NewSwiftCodeHandler tworzy nowy handler
func NewSwiftCodeHandler(repo *repositories.SwiftCodeRepository) *SwiftCodeHandler {
	return &SwiftCodeHandler{repo: repo}
}

// DeleteSwiftCode obsługuje żądanie DELETE /v1/swift-codes/{swift-code}
func (h *SwiftCodeHandler) DeleteSwiftCode(c *gin.Context) {
	swiftCode := strings.ToUpper(strings.TrimSpace(c.Param("swift-code")))

	// Sprawdź, czy kod SWIFT istnieje
	swift, err := h.repo.GetBySwiftCode(swiftCode)
	if err != nil {
		log.Println("Błąd pobierania SWIFT code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Błąd pobierania danych"})
		return
	}
	if swift == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "SWIFT code nie znaleziony"})
		return
	}

	// Jeśli kod to headquarter, usuń powiązania branchy
	if swift.IsHeadquarter {
		if err := h.repo.DetachBranchesFromHeadquarter(swift.ID); err != nil {
			log.Println("Błąd odłączania branchy:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Błąd aktualizacji branchy"})
			return
		}
	}

	// Usuń kod SWIFT
	if err := h.repo.DeleteSwiftCode(swiftCode); err != nil {
		log.Println("Błąd usuwania SWIFT code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Błąd usuwania SWIFT code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code usunięty"})
}
