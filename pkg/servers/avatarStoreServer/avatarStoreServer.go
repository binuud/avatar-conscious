// avatarServer implements the AvatarService for the service from avatar proto file
// This server manages avatars (robots, assistants, etc.)
package avatarStoreServer

import (
	"context"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"

	"github.com/sirupsen/logrus"
)

type avatarGrpcServer struct {
	pb.UnimplementedAvatarStoreServiceServer
	store *avatarStore
}

func NewAvatarStoreGrpcServer() *avatarGrpcServer {
	return &avatarGrpcServer{
		store: NewAvatarStore(),
	}
}

func (s *avatarGrpcServer) List(ctx context.Context, in *pb.AvatarStoreListRequest) (*pb.AvatarStoreListResponse, error) {
	logrus.Infof("ListAvatar called with limit=%d, offset=%d", in.GetLimit(), in.GetOffset())

	avatars, totalCount := s.store.List(in.GetLimit(), in.GetOffset())

	return &pb.AvatarStoreListResponse{
		Avatars:    avatars,
		TotalCount: totalCount,
	}, nil
}

func (s *avatarGrpcServer) Get(ctx context.Context, in *pb.AvatarStoreGetRequest) (*pb.AvatarStoreGetResponse, error) {
	logrus.Infof("GetAvatar called with uuid=%s", in.GetUuid())

	avatar, err := s.store.Get(in.GetUuid())
	if err != nil {
		return nil, err
	}

	return &pb.AvatarStoreGetResponse{
		Avatar: avatar,
	}, nil
}

func (s *avatarGrpcServer) Create(ctx context.Context, in *pb.AvatarStoreCreateRequest) (*pb.AvatarStoreCreateResponse, error) {
	logrus.Infof("CreateAvatar called with name=%s, description=%s, permissions=%v, systemPrompt=%s",
		in.GetName(), in.GetDescription(), in.GetPermissions(), in.GetSystemPrompt())

	avatar := s.store.Create(in.GetName(), in.GetDescription(), in.GetPermissions(), in.GetSystemPrompt())

	return &pb.AvatarStoreCreateResponse{
		Avatar: avatar,
	}, nil
}

func (s *avatarGrpcServer) Update(ctx context.Context, in *pb.AvatarStoreUpdateRequest) (*pb.AvatarStoreUpdateResponse, error) {
	logrus.Infof("UpdateAvatar called with uuid=%s, name=%s, description=%s, permissions=%v, systemPrompt=%s",
		in.GetUuid(), in.GetName(), in.GetDescription(), in.GetPermissions(), in.GetSystemPrompt())

	avatar, err := s.store.Update(in.GetUuid(), in.GetName(), in.GetDescription(), in.GetPermissions(), in.GetSystemPrompt())
	if err != nil {
		return nil, err
	}

	return &pb.AvatarStoreUpdateResponse{
		Avatar: avatar,
	}, nil
}

func (s *avatarGrpcServer) Delete(ctx context.Context, in *pb.AvatarStoreDeleteRequest) (*pb.AvatarStoreDeleteResponse, error) {
	logrus.Infof("DeleteAvatar called with uuid=%s", in.GetUuid())

	avatar, err := s.store.Get(in.GetUuid())
	if err != nil {
		return &pb.AvatarStoreDeleteResponse{
			Message: err.Error(),
		}, err
	}

	// Remove avatar from store
	err = s.store.Delete(avatar.Meta.Uuid)
	if err != nil {
		return &pb.AvatarStoreDeleteResponse{
			Message: err.Error(),
		}, err
	}

	return &pb.AvatarStoreDeleteResponse{
		Message: "Avatar deleted successfully",
	}, nil
}
