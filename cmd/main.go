package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mroczekDNF/swift-api/internal/db"
	"github.com/mroczekDNF/swift-api/internal/routes"
	"github.com/mroczekDNF/swift-api/internal/services"
)

func getEnv(keys ...string) map[string]string {
	values := make(map[string]string)
	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("Missing environment variable: %s", key)
		}
		values[key] = val
	}
	return values
}

func main() {
	// Retrieve environment variables
	envVars := getEnv("DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME")

	// Create DSN and initialize the database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		envVars["DB_HOST"], envVars["DB_USER"], envVars["DB_PASSWORD"], envVars["DB_NAME"], envVars["DB_PORT"])
	db.InitDatabase(dsn)
	defer db.CloseDatabase()
	db.MigrateDatabase()

	// Load data if the table is empty
	if isEmpty, err := db.IsTableEmpty("swift_codes"); err != nil {
		log.Fatalf("Error checking table `swift_codes`: %v", err)
	} else if isEmpty {
		log.Println("Table `swift_codes` is empty. Parsing data...")
		if swiftCodes, err := services.ParseSwiftCodes("data/swift_codes.csv"); err != nil {
			log.Fatalf("Error parsing SWIFT codes: %v", err)
		} else if err := services.SaveSwiftCodesToDatabase(db.DB, swiftCodes); err != nil {
			log.Fatalf("Error saving SWIFT codes to database: %v", err)
		}
		log.Println("Data successfully saved to the database!")
	} else {
		log.Println("Table `swift_codes` contains data. Skipping parsing.")
	}

	// Start the server
	r := routes.SetupRouter(db.DB)
	log.Fatal(r.Run(":8080"))
}
