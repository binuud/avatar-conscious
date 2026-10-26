package datalayer

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"buf.build/go/protoyaml"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// // Pass two type parameters: the base struct T, and P which acts as *T
// func MarshalProtoYAML[T any, P ProtoMessage[T]](msg P) ([]byte, error) {
// 	// msg is explicitly recognized by the compiler as a proto.Message
// 	return protoyaml.Marshal(msg)
// }

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

// SaveYAML write any type 'T' to file in yaml format
func SaveYAML[T any](filePath string, data T) error {

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	// Set custom indent level (e.g., 2 spaces)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Save the YAML data to a file
	// 0644 grants read/write permissions to the owner and read-only to others.
	err := os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil

}

// ReadYAML reads data from a yaml file into a variable of type 'T' and returns it.
func ReadYAML[T any](filePath string) (T, error) {
	var data T // Declares a zero-value variable of type T

	yamlContents, err := os.ReadFile(filePath)
	if err != nil {
		return data, err
	}

	// Unmarshal the YAML data into the struct pointer
	err = yaml.Unmarshal(yamlContents, &data)
	if err != nil {
		return data, err
	}

	return data, nil

}

// SaveProtoYAML write proto.Message type 'T' to file in yaml format
// if we dont use protoyaml marshall, ENUM's are saved as integers
// which is unreadable
func SaveProtoYAML[T proto.Message](filePath string, data T) error {

	yamlData, err := protoyaml.MarshalOptions{
		Indent: 2, // Sets the number of spaces for indentation
	}.Marshal(data)
	if err != nil {
		log.Fatalf("Error occurred during marshaling: %v", err)
	}

	// Save the YAML data to a file
	// 0644 grants read/write permissions to the owner and read-only to others.
	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadProtoYAML reads data from a yaml file into a variable of type 'T' and returns it.
// this might look like a duplicate of read yaml, but simpler implementation
// and less confusing
// binu
func ReadProtoYAML[T proto.Message](filePath string, data T) error {

	// var data T // Declares a zero-value variable of type T

	yamlContents, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("config file not found %v", err)
		return err
	}
	logrus.Infof("Yaml contents %s", yamlContents)
	// Unmarshal the YAML data into the struct pointer
	err = protoyaml.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(yamlContents, data)
	if err != nil {
		return err
	}
	return nil

}
