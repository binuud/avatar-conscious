package avatarServer

import (
	"context"
	"errors"
	"fmt"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	logrus.Infof("Avatar Get called with %s:%s", in.Uuid, in.Name)

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

func (s *avatarGrpcServer) Control(stream grpc.BidiStreamingServer[pb.ControlRequest, pb.ControlResponse]) error {

	logrus.Println("Connect: Client connected to bidirectional stream")

	var a *AvatarWrapper
	var out *pb.ControlResponse

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

		switch in.Control.Action {
		// start the avatar and run to completion
		case pb.AvatarControl_Run:
			{

				a = NewAvatar(in.Control.Meta.Name)

				// initialize the out object
				out = &pb.ControlResponse{
					Control: in.Control,
				}
				out.Control.CreatedAt = utils.GetTimeStamp()
				out.Control.CurrentStep = 0
				out.Control.StepStatuses = a.GetStepStatusesObj()

				if a == nil {
					out.Control.Status = pb.AvatarControl_Fatal
					if err := stream.Send(out); err != nil {
						logrus.Errorf("Connect: Failed to send message: %v", err)
						return err
					}
				}
				a.Start()

				// intimate client on starting the Avatar
				out.Control.Status = pb.AvatarControl_Active
				if err := stream.Send(out); err != nil {
					logrus.Errorf("Connect: Failed to send message: %v", err)
					return err
				}

				for i, step := range a.Avatar.Steps {

					// intimate client on the step completed and its status
					out.Control.CurrentStep = int32(i)
					out.Control.StepStatuses[i].Status = pb.StepStatus_Active
					out.Control.StepStatuses[i].CreatedAt = utils.GetTimeStamp()
					if err := stream.Send(out); err != nil {
						logrus.Errorf("Connect: Failed to send message: %v", err)
						return err
					}

					err := a.RunStep(i, step)

					out.Control.StepStatuses[i].CompletedAt = utils.GetTimeStamp()
					if err != nil {
						out.Control.StepStatuses[i].Status = pb.StepStatus_Fatal
						out.Control.Status = pb.AvatarControl_Fatal
						if err := stream.Send(out); err != nil {
							logrus.Errorf("Connect: Failed to send message: %v", err)
							return err
						}
						return err
					}

					out.Control.StepStatuses[i].Status = pb.StepStatus_Success
					if err := stream.Send(out); err != nil {
						logrus.Errorf("Connect: Failed to send message: %v", err)
						return err
					}

				}
				logrus.Infof("Avatar run completed")
				out.Control.CompletedAt = utils.GetTimeStamp()
				out.Control.Status = pb.AvatarControl_Success
				if err := stream.Send(out); err != nil {
					logrus.Errorf("Connect: Failed to send message: %v", err)
					return err
				}
				return nil
			}
		// run avatar step by step
		// in debug mode
		case pb.AvatarControl_StartDebug:
			{
				a = NewAvatar(in.Control.Meta.Name)

				out = &pb.ControlResponse{
					Control: in.Control,
				}
				out.Control.CreatedAt = utils.GetTimeStamp()
				out.Control.CurrentStep = 0
				out.Control.StepStatuses = a.GetStepStatusesObj()

				if a == nil {
					out.Control.Status = pb.AvatarControl_Fatal
					if err := stream.Send(out); err != nil {
						logrus.Errorf("Connect: Failed to send message: %v", err)
						return err
					}
				}
				a.Start()
				// initialize the out object
				out.Control.Status = pb.AvatarControl_Active
				if err := stream.Send(out); err != nil {
					logrus.Errorf("Connect: Failed to send message: %v", err)
					return err
				}
			}
		case pb.AvatarControl_Stop:
			{
				logrus.Infof("Avatar %s, stopped on user request", a.Meta.Name)
				out.Control.Status = pb.AvatarControl_Failed
				if err := stream.Send(out); err != nil {
					logrus.Errorf("Connect: Failed to send message: %v", err)
					return err
				}

			}
		}

	} // end switch
	return nil

}

func (s *avatarGrpcServer) sendMessage(stream grpc.BidiStreamingServer[pb.ControlRequest, pb.ControlResponse],
	out *pb.ControlResponse) error {

	if stream == nil {
		return fmt.Errorf("stream probably closed")
	}
	if err := stream.Send(out); err != nil {
		logrus.Errorf("Connect: Failed to send message: %v", err)
		return err
	}
	return nil

}
