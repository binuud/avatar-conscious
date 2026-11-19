package memory

import (
	"os"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/sirupsen/logrus"
)

func ListMemories() ([]*pb.MemoryMeta, error) {
	// Read the current directory
	entries, err := os.ReadDir(dataLayer.GetFileAbsolutePath(MEMORY_DIRECTORY, "."))
	if err != nil {
		return nil, err
	}

	memories := make([]*pb.MemoryMeta, 0)
	for _, entry := range entries {
		// Filter out directories if you only want files
		if !entry.IsDir() {
			if entry.Name() == ".DS_Store" {
				continue
			}

			memory, err := LoadMemory(entry.Name())
			if err != nil {
				logrus.Errorf("error cannot load memory %s", memory)
			} else {
				memories = append(memories, memory.Meta)
			}

		}
	}
	return memories, nil
}

func ListMemoryFiles() ([]string, error) {
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
