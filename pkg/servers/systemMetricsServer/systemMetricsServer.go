package systemMetricsServer

import (
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/utils"
)

const defaultInterval = 2 * time.Second

type grpcMetricsServer struct {
	pb.UnimplementedMetricsServer
}

func NewGrpcSystemMetricsServer() *grpcMetricsServer {
	utils.IsGpuMetricsAvailable()
	return &grpcMetricsServer{}
}

// WatchHealth streams periodic health snapshots until the client cancels.
func (s *grpcMetricsServer) GetSystemMetrics(req *pb.GetSystemMetricsRequest, stream pb.Metrics_GetSystemMetricsServer) error {
	interval := defaultInterval
	// interval seconds cannot be less than 2
	if req.GetIntervalSeconds() > 2 {
		interval = time.Duration(req.GetIntervalSeconds()) * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			rep, err := utils.NewSystemMetrics()
			if err != nil {
				return err
			}
			if err := stream.Send(rep); err != nil {
				return err
			}
		}
	}
}
