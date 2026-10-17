package datalayer

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func GetStoreDirectory() string {

	// FILE_STORE_PATH is the folder where data is stored as json files
	FILE_STORE_PATH := os.Getenv("FILE_STORE_PATH")
	if FILE_STORE_PATH == "" {
		logrus.Fatal("env FILE_STORE_PATH not set")
	}
	return FILE_STORE_PATH

}

// GetFileAbsolutePath
// returns a concated string of FILE_STORE_PATH from env
// subDir directory, and the fileName
func GetFileAbsolutePath(subDir string, fileName string) string {

	return fmt.Sprintf("%s/%s/%s", GetStoreDirectory(), subDir, fileName)

}

// SaveJSON saves any type 'T' to the specified file path.
func SaveJSON[T any](filePath string, data T) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print formatting

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	return nil
}

// ReadJSON reads data from a file into a variable of type 'T' and returns it.
func ReadJSON[T any](filePath string) (T, error) {
	var data T // Declares a zero-value variable of type T

	file, err := os.Open(filePath)
	if err != nil {
		return data, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("failed to decode data: %w", err)
	}

	return data, nil
}

// SaveBinary write any type 'T' to file in binary format
func SaveBinary[T any](filePath string, data T) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	return nil
}

// ReadBinary reads data from a file into a variable of type 'T' and returns it.
func ReadBinary[T any](filePath string) (T, error) {
	var data T // Declares a zero-value variable of type T

	file, err := os.Open(filePath)
	if err != nil {
		return data, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("failed to decode data: %w", err)
	}

	return data, nil

}
