package systemMetricsServer

import (
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/sirupsen/logrus"
)

const defaultInterval = 2 * time.Second

type grpcMetricsServer struct {
	pb.UnimplementedMetricsServer
}

func NewGrpcSystemMetricsServer() *grpcMetricsServer {
	return &grpcMetricsServer{}
}

// WatchHealth streams periodic health snapshots until the client cancels.
func (s *grpcMetricsServer) GetSystemMetrics(req *pb.GetSystemMetricsRequest, stream pb.Metrics_GetSystemMetricsServer) error {

	// check if metrics implemented
	_, err := utils.NewSystemMetrics()
	if err != nil {
		logrus.Errorf("Metrics not supported %v", err)
		return err
	} else {
		logrus.Info("Streaming metrics stream request")
	}

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
			logrus.Infof("metrics sse closed from client")
			return stream.Context().Err()
		case <-ticker.C:
			resp, err := utils.NewSystemMetrics()
			if err != nil {
				return err
			}
			// logrus.Infof("metrics response %v", resp)
			if err := stream.Send(resp); err != nil {
				return err
			}
		}
	}
}
