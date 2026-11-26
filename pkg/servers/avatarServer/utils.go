package avatarServer

import (
	"os"

	"github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarStoreServer"
)

func ListAvatars() ([]string, error) {
	// Read the current directory
	entries, err := os.ReadDir(dataLayer.GetFileAbsolutePath(avatarStoreServer.GetAvatarBuildDir(), "."))
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

func ParseAvatarFile(fileName string) {

}
