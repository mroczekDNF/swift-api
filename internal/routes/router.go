package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mroczekDNF/swift-api/internal/handlers"
)

// SetupRouter definiuje endpointy API
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Endpoint testowy
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })
	router.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCodeDetails)

	return router
}
