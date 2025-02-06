package handlers

import "github.com/mroczekDNF/swift-api/internal/repositories"

// SwiftCodeHandler handles operations on SWIFT codes.
type SwiftCodeHandler struct {
	repo repositories.SwiftCodeRepositoryInterface
}

// NewSwiftCodeHandler creates a new handler.
func NewSwiftCodeHandler(repo repositories.SwiftCodeRepositoryInterface) *SwiftCodeHandler {
	return &SwiftCodeHandler{repo: repo}
}
