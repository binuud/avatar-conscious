package utils

import (
	"fmt"
	"path/filepath"
	"strings"

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

func GetUploadsFolder() string {
	return GetEnv("FILE_STORE_PATH", ".") + "/uploads/"
}

func GetGeneratedFilesFolder() string {
	return GetEnv("FILE_STORE_PATH", ".") + "/gen/"
}
