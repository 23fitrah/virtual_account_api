package config

import (
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logEntry := map[string]interface{}{
			"level":     "error",
			"timestamp": time.Now().Format(time.RFC3339),
			"message":   ".env file not found or failed to load",
			"error":     err.Error(),
		}

		jsonBytes, _ := sonic.Marshal(logEntry)

		fmt.Println(string(jsonBytes))
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
