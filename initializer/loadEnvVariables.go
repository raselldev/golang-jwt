package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

// The function `LoadEnvVariables` loads environment variables from a .env file in Go.
func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
