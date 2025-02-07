package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DeleteSwiftCode handles DELETE /v1/swift-codes/{swift-code} requests.
func (h *SwiftCodeHandler) DeleteSwiftCode(c *gin.Context) {
	swiftCode := strings.ToUpper(strings.TrimSpace(c.Param("swift-code")))

	swift, err := h.repo.GetBySwiftCode(swiftCode)
	if err != nil {
		log.Printf("Error retrieving SWIFT code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving SWIFT code"})
		return
	}
	if swift == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "SWIFT code not found"})
		return
	}

	if swift.IsHeadquarter {
		if err := h.repo.DetachBranchesFromHeadquarter(swift.ID); err != nil {
			log.Printf("Error detaching branches: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error detaching branches"})
			return
		}
	}

	if err := h.repo.DeleteSwiftCode(swiftCode); err != nil {
		log.Printf("Error deleting SWIFT code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting SWIFT code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully"})
}
