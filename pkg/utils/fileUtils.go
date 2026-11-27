package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// save the image base 64 into uploads folder
// provide a unique name to the fileName, use UUID always
func SaveImageDataToFile(fileName string, rawBase64 string) error {

	// Fix Corruption: Locate and strip the "data:image/png;base64," prefix if it exists
	if idx := strings.Index(rawBase64, ","); idx != -1 {
		rawBase64 = rawBase64[idx+1:]
	}

	// 2. Decode the clean, sanitized base64 payload string
	decodedBytes, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		return status.Error(codes.InvalidArgument, "saving image to file, cannot decode file: "+err.Error())
	}

	if len(decodedBytes) == 0 {
		return status.Error(codes.InvalidArgument, "saving image to file, file empty")
	}

	destinationPath := GetEnv("FILE_STORE_PATH", ".") + "/uploads/" + fileName
	logrus.Infof("Uplaoded file to be saved under -> %s", destinationPath)

	err = os.WriteFile(destinationPath, decodedBytes, 0644)
	if err != nil {
		return status.Error(codes.InvalidArgument, "cannot write file to server: "+err.Error())
	}
	return nil
}

func FileExists(filename string) bool {
	// logrus.Infof("Checking file exists %f", filename)
	_, err := os.Stat(filename)
	if err == nil {
		return true // File exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // File explicitly does not exist
	}
	// Any other error means the file might exist, but we couldn't stat it
	// (e.g., permission denied)
	return false
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func CopyFile(srcPath, dstPath string) error {
	// Open the source file for reading
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Ensure the destination folder exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create (or overwrite) the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Stream the file contents from source to destination
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	// Flush writes to stable storage
	return dstFile.Sync()
}

func CopyDir(srcDir, dstDir string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("Error creating destination:", err)
	}

	// Copy the contents of srcDir into dstDir
	err := os.CopyFS(dstDir, os.DirFS(srcDir))
	if err != nil {
		return fmt.Errorf("Error copying directory:", err)
	}
	return nil
}
