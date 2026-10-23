package avatarServer

import (
	"os"

	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
)

func ListAvatars() ([]string, error) {
	// Read the current directory
	entries, err := os.ReadDir(datalayer.GetFileAbsolutePath(AVATAR_DIRECTORY, "."))
	if err != nil {
		return nil, err
	}

	items := make([]string, 0)
	for _, entry := range entries {
		// Filter out directories if you only want files
		if !entry.IsDir() {
			items = append(items, entry.Name())
		}
	}
	return items, nil
}
