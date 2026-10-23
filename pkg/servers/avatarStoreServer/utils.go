package avatarStoreServer

import (
	"log"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
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
			avatar, err := ReadAvatarFromStore(entry.Name())
			if err == nil {
				items = append(items, avatar)
			}
		}
	}
	return items, nil
}

func ReadAvatarFromStore(fileName string) (*pb.Avatar, error) {

	if !strings.HasSuffix(fileName, ".md") {
		fileName = fileName + ".md"
	}

	filePath := datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// fileBytes, err := os.ReadFile(filePath)
	// if err != nil {
	// 	logrus.Errorf("Error reading %s: %v\n", fileName, err)
	// 	return nil, err
	// }

	var meta pb.AvatarMeta

	// 3. Parse the frontmatter and extract the remaining content body
	// frontmatter.Parse reads from any io.Reader, populates the struct,
	// and returns the rest of the file content as a byte slice.
	content, err := frontmatter.Parse(file, &meta)
	if err != nil {
		log.Fatalf("Failed to parse frontmatter: %v", err)
	}

	return &pb.Avatar{
		Meta:    &meta,
		Content: string(content),
	}, nil
}
