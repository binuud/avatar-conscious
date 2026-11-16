package memory

import (
	"os"

	"github.com/binuud/avatar-conscious/pkg/dataLayer"
)

func ListMemory() ([]string, error) {
	// Read the current directory
	entries, err := os.ReadDir(dataLayer.GetFileAbsolutePath(MEMORY_DIRECTORY, "."))
	if err != nil {
		return nil, err
	}

	memories := make([]string, 0)
	for _, entry := range entries {
		// Filter out directories if you only want files
		if !entry.IsDir() {
			memories = append(memories, entry.Name())
		}
	}
	return memories, nil
}

func DeleteMemory(name string) error {
	return os.Remove(dataLayer.GetFileAbsolutePath(MEMORY_DIRECTORY, name))
}
