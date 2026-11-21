package settingsServer

import (
	"context"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/sirupsen/logrus"
)

type grpcSettingsServer struct {
	pb.UnimplementedSettingsServer
}

func NewGrpcSettingsServer() *grpcSettingsServer {
	return &grpcSettingsServer{}
}

func (s *grpcSettingsServer) List(ctx context.Context, in *pb.SettingsListRequest) (*pb.SettingsListResponse, error) {

	logrus.Debug("Settings Server List")
	prefillFields(settings.Global)
	settings.Global.Save()
	// time.Sleep(5 * time.Second) // Pauses for 2 seconds
	return &pb.SettingsListResponse{
		Groups: settings.Global.GetAllByGroup(),
	}, nil
}

func (s *grpcSettingsServer) Save(ctx context.Context, in *pb.SettingsSaveRequest) (*pb.SettingsSaveResponse, error) {

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
