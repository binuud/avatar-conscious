package avatarServer

import (
	"context"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/sirupsen/logrus"
)

type avatarGrpcServer struct {
	pb.UnimplementedAvatarServiceServer
	avatar *pb.Avatar
}

func NewAvatarGrpcServer() *avatarGrpcServer {
	return &avatarGrpcServer{
		avatar: &pb.Avatar{},
	}
}

func (s *avatarGrpcServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("Avatar start called with %s:%s", in.Uuid, in.Name)

	a := NewAvatar(in.Name)
	return &pb.GetResponse{
		Avatar: a.Avatar,
	}, nil
}

func (s *avatarGrpcServer) Start(ctx context.Context, in *pb.StartRequest) (*pb.StartResponse, error) {
	logrus.Infof("Avatar start called with %s:%s", in.Uuid, in.Name)

	a := NewAvatar(in.Name)
	return &pb.StartResponse{
		Avatar: a.Avatar,
	}, nil
}

func (s *avatarGrpcServer) Stop(ctx context.Context, in *pb.StopRequest) (*pb.StopResponse, error) {
	logrus.Infof("Avatar stop called with %s", in.Uuid)

	return &pb.StopResponse{}, nil
}

func (s *avatarGrpcServer) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	logrus.Infof("Avatar connect called with %s", in.Uuid)

	return &pb.ConnectResponse{}, nil
}
