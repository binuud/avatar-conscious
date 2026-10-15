package utils

import (
	"os"
	"time"
	"uuid"

	"github.com/binuud/conscious/gen/go/v1/conscious"
)

func GenerateUUID() string {
	return uuid.New().String()
}

// getEnv retrieves an environment variable or returns a fallback default value
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetTimeStamp() *conscious.Timestamp {

	return &conscious.Timestamp{
		Seconds: time.Now().UnixMilli(),
		Nanos:   0,
	}

}
