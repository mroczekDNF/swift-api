package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
)

// SetupRouter definiuje endpointy API
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCodeDetails)
	router.GET("/v1/swift-codes/country/:country", handlers.GetSwiftCodesByCountry)

	return router
}
