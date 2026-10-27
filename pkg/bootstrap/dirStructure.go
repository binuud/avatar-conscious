package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// CreateStorageDirStructure
// create the needed directory structure under storage path
func CreateStorageDirStructure() {

	FILE_STORE_PATH := os.Getenv("FILE_STORE_PATH")
	if FILE_STORE_PATH == "" {
		logrus.Fatal("env FILE_STORE_PATH not set")
	}
	// Define the directory structure you need on startup
	dirsToCreate := []string{
		filepath.Join(FILE_STORE_PATH, "avatars"),
		filepath.Join(FILE_STORE_PATH, "avatarStore"),
		filepath.Join(FILE_STORE_PATH, "configs"),
		filepath.Join(FILE_STORE_PATH, "memories"),
		filepath.Join(FILE_STORE_PATH, "mockServer"),
		filepath.Join(FILE_STORE_PATH, "sampleImages"),
		filepath.Join(FILE_STORE_PATH, "static", "avatarStore"),
		filepath.Join(FILE_STORE_PATH, "static", "gen", "images"),
		filepath.Join(FILE_STORE_PATH, "static", "gen", "tts"),
		filepath.Join(FILE_STORE_PATH, "uploads"),
	}

	// Initialize the file structure
	logrus.Println("Initializing avatar storage directory structure...")
	if err := initStorageStructure(dirsToCreate); err != nil {
		logrus.Fatalf("Failed to initialize storage directories: %v", err)
	}

	logrus.Println("Successfully verified/created directory structure. Proceeding with application start...")

}

// initStorageStructure loops through target paths and securely creates them
func initStorageStructure(directories []string) error {
	for _, dir := range directories {
		// os.ModePerm (0777) applies broad read/write/execute permissions,
		// which are modified by the system's umask for security.
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
