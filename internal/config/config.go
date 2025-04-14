package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {

	// load .env file with godotenv might be problematic with docker, bcuz docker use real enviroment variable
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
