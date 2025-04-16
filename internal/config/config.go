package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error Loading .env: %v", err)
	}
	return os.Getenv(key)
}
