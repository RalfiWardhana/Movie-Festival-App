package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// InitEnvironment loads environment variables from .env file
func InitEnvironment() {
	// Load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Validate if JWT_SECRET is set
	if len(os.Getenv("JWT_SECRET")) == 0 {
		log.Fatal("JWT_SECRET is not set in the environment variables")
	}
}
