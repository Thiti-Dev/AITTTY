package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig -> Load a config value from the given config key
func LoadConfig(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Failed loading env file")
	}
	return os.Getenv(key)
}