package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/binuud/avatar-conscious/pkg/settings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- Security Logic ---

// Safe extensions allowlist. Modify this list based on your needs.
var _safeExtensions = map[string]bool{
	"txt":  true,
	"pdf":  true,
	"png":  true,
	"jpg":  true,
	"jpeg": true,
	"gif":  true,
	"mp3":  true,
	"mp4":  true,
	"doc":  true,
	"docx": true,
	"xlsx": true,
	"csv":  true,
	"json": true,
	"xml":  true,
}

// Allowed file types that are explicitly dangerous (for double-checking)
var _dangerousExtensions = map[string]bool{
	"exe": true,
	"sh":  true,
	"bat": true,
	"cmd": true,
	"ps1": true, // PowerShell
	"vbs": true, // Visual Basic
	"msi": true,
	"scr": true,
	"hta": true,
	"lnk": true,
	"jar": true,
	"dll": true,
	"so":  true,
}

// Check if a filename is safe based on extension
func IsFilenameSafe(fileName string) (bool, error) {

	if fileName == "" {
		return false, status.Error(codes.InvalidArgument, "filename cannot be empty")
	}

	if !strings.Contains(fileName, ".") {
		return false, status.Error(codes.InvalidArgument, "filename must have an extension")
	}

	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)
	// Remove the dot
	ext = strings.TrimPrefix(ext, ".")

	if _dangerousExtensions[ext] {
		return false, fmt.Errorf("unsafe extension: %s", ext)
	}

	if !_safeExtensions[ext] {
		return false, fmt.Errorf("unsafe extension: %s", ext)
	}

	return true, nil
}

// GetFileAbsolutePath
// if relativePath starts with /static, return the absolute path of static subfolder + relativefilepath
// else check if file is in uploads folder, and send path, else return error
func GetFileAbsolutePath(relativePath string) (string, error) {
	if strings.HasPrefix(relativePath, "/static") {
		filePath := settings.Global.GetValue("SERVABLE_DIRECTORY") + relativePath[7:] // pop the string /static from the relative path
		if FileExists(filePath) {
			return filePath, nil
		}
	}
	filePath := GetUploadsFolder() + relativePath
	if FileExists(filePath) {
		return filePath, nil
	}
	return "", fmt.Errorf("filepath %s, not found or not accessible", relativePath)
}

func GetUploadsFolder() string {
	return settings.Global.GetValue("FILE_STORE_PATH") + "/uploads/"
}

func GetGeneratedFilesFolder() string {
	return settings.Global.GetValue("FILE_STORE_PATH") + "/gen/"
}

func GetSamplesFolder() string {
	return settings.Global.GetValue("FILE_STORE_PATH") + "/sampleImages/"
}
