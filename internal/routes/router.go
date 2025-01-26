package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
)

// SetupRouter definiuje endpointy API
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/v1/swift-codes/:swiftCode", handlers.GetSwiftCodeDetails)
	router.GET("/v1/swift-codes/country/:countryISO2", handlers.GetSwiftCodesByCountry)
	router.POST("/v1/swift-codes", handlers.AddSwiftCode)

	return router
}
