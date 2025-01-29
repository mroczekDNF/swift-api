package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
	"github.com/mroczekDNF/swift-api/internal/repositories"
)

// SetupRouter definiuje endpointy API
func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	repo := repositories.NewSwiftCodeRepository(db)
	handler := handlers.NewSwiftCodeHandler(repo)

	router.GET("/v1/swift-codes/:swiftCode", handler.GetSwiftCodeDetails)
	router.GET("/v1/swift-codes/country/:countryISO2", handler.GetSwiftCodesByCountry)
	router.POST("/v1/swift-codes", handler.AddSwiftCode)
	router.DELETE("/v1/swift-codes/:swift-code", handler.DeleteSwiftCode)

	return router
}
