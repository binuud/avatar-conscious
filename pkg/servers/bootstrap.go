package servers

import (
	"fmt"
	"log"
	"net"

	pbAvatar "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	pbConscious "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	pbSettings "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarServer"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarStoreServer"
	"github.com/binuud/avatar-conscious/pkg/servers/consciousServer"
	"github.com/binuud/avatar-conscious/pkg/servers/settingsServer"
	"github.com/binuud/avatar-conscious/pkg/servers/systemMetricsServer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Bootstrap_GRPC_Services(server_grpc_port int) {

	insecureConn, err := net.Listen("tcp", fmt.Sprintf(":%d", server_grpc_port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	insecureServer := grpc.NewServer()

	// GRPC Conscious Server
	pbConscious.RegisterConsciousServer(insecureServer, consciousServer.NewGrpcConsciousServer())

	// GRPC Avatar Server
	pbAvatar.RegisterAvatarServiceServer(insecureServer, avatarServer.NewAvatarGrpcServer())

	// GRPC Avatar Store Server
	pbAvatar.RegisterAvatarStoreServiceServer(insecureServer, avatarStoreServer.NewAvatarStoreGrpcServer())

	// GRPC SystemMetrics Server
	pbConscious.RegisterMetricsServer(insecureServer, systemMetricsServer.NewGrpcSystemMetricsServer())

	// GRPC Settings Server
	pbSettings.RegisterSettingsServer(insecureServer, settingsServer.NewSettingsGrpcServer())

	// Register reflection service on gRPC server.
	reflection.Register(insecureServer)

	logrus.Infof("GRPC server starting  :(%d)", server_grpc_port)
	logrus.Fatalln(insecureServer.Serve(insecureConn))

}
