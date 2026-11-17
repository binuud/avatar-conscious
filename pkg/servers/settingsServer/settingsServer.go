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

	logrus.Debug("Settings Server List")

	return &pb.SettingsListResponse{
		Groups: settings.Global.GetAllByGroup(),
	}, nil
}

func (s *settingsGrpcServer) Save(ctx context.Context, in *pb.SettingsSaveRequest) (*pb.SettingsSaveResponse, error) {

	logrus.Debugf("Settings Server Save")

	for _, entry := range in.Configs.Configs {
		settings.Global.Update(entry.Name, entry.Value)
	}
	settings.Global.Save()
	logrus.Infof("Settings Saved by group %s", in.GroupKey)

	return &pb.SettingsSaveResponse{
		Configs:  in.Configs,
		GroupKey: in.GroupKey,
		Success:  true,
	}, nil
}
