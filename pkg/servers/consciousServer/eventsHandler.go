package consciousServer

import (
	"errors"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

func (s *grpcConsciousServer) handleSensorEvents(stream grpc.BidiStreamingServer[pb.ContextMesssage, pb.ContextMesssage]) error {
	var out *pb.ContextMesssage
	for {
		// Read message from the client stream
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Client safely closed the sending side of the stream
			logrus.Errorf("Connect: Client closed the connection stream (EOF)")
			return nil
		}
		if err != nil {
			logrus.Errorf("Failed to receive message: %v", err)
			return err
		}
		logrus.Infof("Received action from client %v, Step No: %d ", in.Name)

		out = &pb.ContextMesssage{
			Name: in.Name,
		}
		if err := stream.Send(out); err != nil {
			logrus.Errorf("Connect: Failed to send message: %v", err)
			return err
		}

	}

}
