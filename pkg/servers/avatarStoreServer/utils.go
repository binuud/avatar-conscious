package avatarStoreServer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"buf.build/go/protoyaml"
	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/utils"
)

func ListAvatarsFromStore() ([]*pb.Avatar, error) {
	// Read the current directory
	entries, err := os.ReadDir(datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, "."))
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Avatar, 0)
	for _, entry := range entries {
		// Filter out directories if you only want files
		if !entry.IsDir() {
			// Filter for .md extension (filepath.Ext returns strings like ".md")
			if filepath.Ext(entry.Name()) == ".md" {
				// Construct the full path or just use entry.Name() depending on your needs
				avatar, err := ReadAvatarFromStore(entry.Name())
				if err == nil {
					items = append(items, avatar)
				}
			}

		}
	}
	return items, nil
}

// ReadAvatarFromStore()
// Reads the Avatar md file
// Parses out the yaml part and other md file content
// Unmarshalls yaml into avatar struct
// DO NOT CHANGE THIS
func ReadAvatarFromStore(fileName string) (*pb.Avatar, error) {

	if !strings.HasSuffix(fileName, ".md") {
		fileName = fileName + ".md"
	}

	filePath := datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, fileName)

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var meta pb.AvatarMeta

	yaml, content, ok := utils.ExtractFrontMatter(string(fileContents))
	if !ok {
		return nil, fmt.Errorf("Failed to parse frontmatter, please check if yaml is delimitted by '---'")
	}

	// to get around string to enum convertion
	err = protoyaml.Unmarshal([]byte(yaml), &meta)
	if err != nil {
		return nil, err
	}

	return &pb.Avatar{
		Meta:    &meta,
		Content: string(content),
	}, nil
}
