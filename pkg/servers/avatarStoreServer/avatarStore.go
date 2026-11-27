package avatarStoreServer

import (
	"cmp"
	"fmt"
	"path/filepath"
	"slices"
	"sync"
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type avatarStore struct {
	sync.RWMutex
	avatars map[string]*pb.Avatar
}

func NewAvatarStore() *avatarStore {

	entries, err := ListAvatarsFromStore()
	if err != nil {
		logrus.Errorf("error: NewAvatarStore cannot list avatars %v", err)
		return nil
	}

	avatars := make(map[string]*pb.Avatar)

	for _, entry := range entries {
		avatars[entry.Meta.Uuid] = entry
		logrus.Infof("avatar store %s", entry.Meta.Name)
	}

	return &avatarStore{
		avatars: avatars,
	}
}

func (s *avatarStore) Load() error {
	if err := s.saveToDisk(); err != nil {
		return err
	}
	return nil
}

func (s *avatarStore) saveToDisk() error {
	return nil // TODO: Implement persistence
}

func (s *avatarStore) Save(uuid string) error {
	s.RLock()
	_, exists := s.avatars[uuid]
	s.RUnlock()

	if !exists {
		return status.Error(codes.NotFound, "avatar not found")
	}

	return s.saveToDisk()
}

func (s *avatarStore) Get(uuid string) (*pb.Avatar, error) {
	s.RLock()
	avatar, exists := s.avatars[uuid]
	s.RUnlock()

	if !exists {
		return nil, status.Error(codes.NotFound, "avatar not found")
	}

	return avatar, nil
}

func (s *avatarStore) List(limit, offset int32) ([]*pb.Avatar, int32) {
	s.RLock()
	defer s.RUnlock()

	avatars := make([]*pb.Avatar, 0)

	count := 0

	// Simple iteration (implement proper sorting/filtering as needed)
	for _, avatar := range s.avatars {
		if count >= int(limit) {
			break
		}
		// if count >= int(offset) {
		// 	break
		// }
		avatars = append(avatars, avatar)
		count++
	}

	// Sort by Age ascending
	slices.SortFunc(avatars, func(a, b *pb.Avatar) int {
		return cmp.Compare(a.Meta.Name, b.Meta.Name)
	})

	return avatars, int32(len(s.avatars))
}

func (s *avatarStore) Create(name, description string, permissions pb.Permission, systemPrompt string) *pb.Avatar {
	s.Lock()
	defer s.Unlock()

	uuid := fmt.Sprintf("avatar-%d", time.Now().UnixNano())

	avatar := &pb.Avatar{
		Meta: &pb.AvatarMeta{
			Uuid:        utils.GenerateUUID(),
			Name:        name,
			Description: description,
		},
		Permissions:  permissions,
		SystemPrompt: systemPrompt,
	}

	s.avatars[uuid] = avatar
	s.Save(uuid)

	return avatar
}

func (s *avatarStore) Update(uuid string, name, description string, permissions pb.Permission, systemPrompt string) (*pb.Avatar, error) {
	s.Lock()
	defer s.Unlock()

	avatar, exists := s.avatars[uuid]
	if !exists {
		return nil, status.Error(codes.NotFound, "avatar not found")
	}

	avatar.Meta.Name = name
	avatar.Meta.Description = description
	avatar.Permissions = permissions
	avatar.SystemPrompt = systemPrompt

	s.avatars[uuid] = avatar
	s.Save(uuid)

	return avatar, nil
}

func (s *avatarStore) Delete(uuid string) error {
	s.Lock()
	defer s.Unlock()

	_, exists := s.avatars[uuid]
	if !exists {
		return status.Error(codes.NotFound, "avatar not found")
	}

	delete(s.avatars, uuid)
	s.Save(uuid)

	return nil
}

// Sync
// copies thumbnails and other static files from avatar-store
func (s *avatarStore) Sync() error {
	logrus.Infof("Sync %d", len(s.avatars))
	FILE_STORE_PATH := settings.Global.GetValue("FILE_STORE_PATH")
	// Simple iteration (implement proper sorting/filtering as needed)
	for _, avatar := range s.avatars {
		// copy thumb nail to static folder
		avatarThumbNail := filepath.Base(avatar.Meta.Image)
		avatarDir := GetAvatarDir(avatar.Meta.Name)
		srcThumbnailFile := filepath.Join(avatarDir, avatarThumbNail)
		destThumbnailFile := filepath.Join(FILE_STORE_PATH, "static", "avatarStore", avatar.Meta.Name, avatar.Meta.Image)
		logrus.Infof("Copying thumbnail %s -> %s", srcThumbnailFile, destThumbnailFile)
		utils.CopyFile(srcThumbnailFile, destThumbnailFile)

		// now check if asset folder is present
		assetsFolder := filepath.Join(avatarDir, "assets")

		if utils.DirExists(assetsFolder) {
			destAssetsDir := filepath.Join(FILE_STORE_PATH, "static", "avatarStore", avatar.Meta.Name, "assets")
			logrus.Debugf("Copying assets from %s->%s", assetsFolder, destAssetsDir)
			utils.CopyDir(assetsFolder, destAssetsDir)
		}
	}

	return nil
}
