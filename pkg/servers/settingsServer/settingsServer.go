package settingsServer

import (
	"context"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/sirupsen/logrus"
)

type settingsGrpcServer struct {
	pb.UnimplementedSettingsServer
}

func NewSettingsGrpcServer() *settingsGrpcServer {
	return &settingsGrpcServer{}
}

func (s *settingsGrpcServer) List(ctx context.Context, in *pb.SettingsListRequest) (*pb.SettingsListResponse, error) {

	logrus.Infof("Settings Server List")

	return &pb.SettingsListResponse{
		Groups: settings.Global.GetAllByGroup(),
	}, nil
}
