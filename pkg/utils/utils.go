package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
	"uuid"

	"github.com/binuud/avatar-conscious/gen/go/v1/common"
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

func GetTimeStamp() *common.Timestamp {

	return &common.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   0,
	}

}

func GetBaseUrl(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	baseURL := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	}

	return baseURL.String(), nil
}

func GetDataFromUrl(rawURL string) (string, error) {
	// Make the HTTP GET request
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error Failed to fetch URL: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error Unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error Failed to read response body: %v", err)
	}

	return string(body), nil
}
