package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func GetPrettyJson(rawJSON string) (string, error) {
	var prettyJSON bytes.Buffer
	// Parameters: destination buffer, source bytes, prefix string, indent string
	err := json.Indent(&prettyJSON, []byte(rawJSON), "", "  ")
	if err != nil {
		return "", fmt.Errorf("Invalid JSON: %v", err)
	}
	return prettyJSON.String(), nil
}
