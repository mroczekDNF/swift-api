package handlers

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter definiuje endpointy API
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Endpoint testowy
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return router
}
