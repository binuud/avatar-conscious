package avatarStoreServer

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AVATAR_DIRECTORY       = "avatars"
	AVATAR_STORE_DIRECTORY = "avatarStore"
)

type avatarStore struct {
	sync.RWMutex
	avatars map[string]*pb.Avatar
}

func NewAvatarStore() *avatarStore {

	entries, err := ListAvatarsFromStore()
	if err != nil {
		return nil
	}

	avatars := make(map[string]*pb.Avatar)

	for _, entry := range entries {
		avatars[entry.Meta.Uuid] = entry
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
