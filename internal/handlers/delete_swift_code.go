package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/internal/repositories"
)

// SwiftCodeHandler handles operations on SWIFT codes.
type SwiftCodeHandler struct {
	repo *repositories.SwiftCodeRepository
}

// NewSwiftCodeHandler creates a new handler.
func NewSwiftCodeHandler(repo *repositories.SwiftCodeRepository) *SwiftCodeHandler {
	return &SwiftCodeHandler{repo: repo}
}

// DeleteSwiftCode handles DELETE /v1/swift-codes/{swift-code} requests.
func (h *SwiftCodeHandler) DeleteSwiftCode(c *gin.Context) {
	swiftCode := h.getNormalizedSwiftCode(c)

	swift, err := h.fetchSwiftCode(swiftCode)
	if err != nil {
		h.logErrorAndRespond(c, "Error retrieving SWIFT code", err, http.StatusInternalServerError)
		return
	}
	if swift == nil {
		h.respondNotFound(c, "SWIFT code not found")
		return
	}

	if err := h.detachBranchesIfHeadquarter(swift); err != nil {
		h.logErrorAndRespond(c, "Error detaching branches", err, http.StatusInternalServerError)
		return
	}

	if err := h.deleteSwiftCode(swiftCode); err != nil {
		h.logErrorAndRespond(c, "Error deleting SWIFT code", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully"})
}

// getNormalizedSwiftCode retrieves and normalizes the swift code from the URL parameter.
func (h *SwiftCodeHandler) getNormalizedSwiftCode(c *gin.Context) string {
	return strings.ToUpper(strings.TrimSpace(c.Param("swift-code")))
}

// fetchSwiftCode retrieves the SWIFT code details from the repository.
func (h *SwiftCodeHandler) fetchSwiftCode(swiftCode string) (*models.SwiftCode, error) {
	return h.repo.GetBySwiftCode(swiftCode)
}

// detachBranchesIfHeadquarter detaches branches from a record if it is a headquarter.
func (h *SwiftCodeHandler) detachBranchesIfHeadquarter(swift *models.SwiftCode) error {
	if swift.IsHeadquarter {
		return h.repo.DetachBranchesFromHeadquarter(swift.ID)
	}
	return nil
}

// deleteSwiftCode deletes the SWIFT code using the repository.
func (h *SwiftCodeHandler) deleteSwiftCode(swiftCode string) error {
	return h.repo.DeleteSwiftCode(swiftCode)
}

// logErrorAndRespond logs the error and sends a JSON response.
func (h *SwiftCodeHandler) logErrorAndRespond(c *gin.Context, message string, err error, statusCode int) {
	log.Printf("%s: %v", message, err)
	c.JSON(statusCode, gin.H{"message": message})
}

// respondNotFound sends a not found JSON response.
func (h *SwiftCodeHandler) respondNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"message": message})
}
