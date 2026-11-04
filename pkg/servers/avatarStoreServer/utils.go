package avatarStoreServer

// -- avatar-name.md (with steps and front matter)
// -- avatar-name.ui.yaml (ui forms and pages)
// -- avatar-name.png (the store image)

import (
	"fmt"
	"os"

	"buf.build/go/protoyaml"
	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/sirupsen/logrus"
)

// ListAvatarsFromStore()
// Folder structire
// avatar-name (folder)
// -- avatar-name.md (with steps and front matter)
// -- avatar-name.ui.yaml (ui forms and pages)
// -- avatar-name.png (the store image)
func ListAvatarsFromStore() ([]*pb.Avatar, error) {
	// Read the current directory
	entries, err := os.ReadDir(datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, "."))
	if err != nil {
		return nil, err
	}

	items := make([]*pb.Avatar, 0)
	for _, entry := range entries {

		// all avatars are expected to be directories
		if entry.IsDir() {
			// Filter for .md extension (filepath.Ext returns strings like ".md")
			logrus.Infof("Analysing avatar file %s", entry.Name())

			// Construct the full path or just use entry.Name() depending on your needs
			avatar, err := ReadAvatarFromStore(entry.Name())
			if err == nil {
				items = append(items, avatar)
			} else {
				logrus.Errorf("avatar file %s.md file cannot be read %v", entry.Name(), err)
			}
		}
	}

	return items, nil
}

// ReadAvatarFromStore()
// Reads the Avatar file
// .md is appended to avatar name to get md file
// Parses out the yaml part and other md file content
// Unmarshalls yaml into avatar struct
// avatar-name (folder)

// DO NOT CHANGE THIS
func ReadAvatarFromStore(avatarName string) (*pb.Avatar, error) {

	// the directory containing all avatar files
	avatarDir := datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, avatarName)

	avatarFile := fmt.Sprintf("%s/%s.md", avatarDir, avatarName)
	avatarFormFile := fmt.Sprintf("%s/%s-ui.yaml", avatarDir, avatarName)

	// read the front matter, and create the avatar object
	fileContents, err := os.ReadFile(avatarFile)
	if err != nil {
		logrus.Fatalf("Failed to read file: %s", err)
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

	uiCollection := &pb.UiCollection{}
	err = datalayer.ReadProtoYAML(avatarFormFile, uiCollection)
	if err != nil {
		logrus.Errorf("Failed to read file: %s", err)
	}

	return &pb.Avatar{
		Meta:         &meta,
		Content:      string(content),
		UiCollection: uiCollection,
	}, nil
}
