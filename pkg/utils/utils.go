package utils

import "os"

// getEnv retrieves an environment variable or returns a fallback default value
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
