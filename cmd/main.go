package main

import (
	"log"

	"github.com/mroczekDNF/swift-api/internal/handlers"
)

func main() {
	router := handlers.SetupRouter()

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
