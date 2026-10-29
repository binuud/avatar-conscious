package avatarServer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"uuid"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	pbConscious "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type avatarGrpcServer struct {
	pb.UnimplementedAvatarServiceServer
	avatar *AvatarWrapper
}

func NewAvatarGrpcServer() *avatarGrpcServer {
	return &avatarGrpcServer{
		avatar: &AvatarWrapper{},
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
		logrus.Infof("Received action from client %v, %v", in.Control.Action, in)

		switch in.Control.Action {
		// initialize the avatar session
		case pb.AvatarControl_INIT:
			{
				err := s.handleAvatarInit(in, stream)
				if err != nil {
					return err
				}
			}
		// start the avatar and run to completion
		// break if ui intervention is needed
		case pb.AvatarControl_RUN:
			{
				err := s.handleAvatarRun(in, stream)
				if err != nil {
					return err
				}
			}

		// run avatar step by step
		// in debug mode
		case pb.AvatarControl_START_DEBUG:
			{

				out = &pb.ControlResponse{
					Control: in.Control,
				}
				out.Control.CreatedAt = utils.GetTimeStamp()
				out.Control.CurrentStep = 0

				if s.avatar == nil {
					out.Control.Status = pb.AvatarStatus_STATUS_ACTIVE
					if err := stream.Send(out); err != nil {
						logrus.Errorf("Connect: Failed to send message: %v", err)
						return err
					}
				}

				out.Control.StepStatuses = s.avatar.GetStepStatusesObj()

				s.avatar.Start()
				// initialize the out object
				out.Control.Status = pb.AvatarStatus_STATUS_ACTIVE
				if err := stream.Send(out); err != nil {
					logrus.Errorf("Connect: Failed to send message: %v", err)
					return err
				}
			}
		case pb.AvatarControl_STOP:
			{
				logrus.Infof("Avatar %s, stopped on user request", s.avatar.Meta.Name)
				out.Control.Status = pb.AvatarStatus_STATUS_FAILED
				if err := stream.Send(out); err != nil {
					logrus.Errorf("Connect: Failed to send message: %v", err)
					return err
				}

			}
		}

	} // end switch

}

func (s *avatarGrpcServer) handleAvatarInit(in *pb.ControlRequest,
	stream grpc.BidiStreamingServer[pb.ControlRequest, pb.ControlResponse]) error {
	s.avatar = NewAvatar(in.Control.Meta.Name)

	// initialize memory
	s.avatar.Start()
	myMsg := proto.Clone(in.Control)
	// initialize the out object
	out := &pb.ControlResponse{
		Control: myMsg.(*pb.AvatarControl),
	}
	out.Control.CurrentStep = 0
	if s.avatar == nil {
		out.Control.Status = pb.AvatarStatus_STATUS_FATAL
	} else {
		// intimate client on initializing the Avatar
		out.Control.Status = pb.AvatarStatus_STATUS_INIT
	}

	if err := stream.Send(out); err != nil {
		logrus.Errorf("Connect: Failed to send message: %v", err)
		return err
	}
	return nil
}

func (s *avatarGrpcServer) handleAvatarRun(in *pb.ControlRequest,
	stream grpc.BidiStreamingServer[pb.ControlRequest, pb.ControlResponse]) error {

	// clone the in data, because we dont want to loose it
	// when out msg is changed
	myMsg := proto.Clone(in.Control)
	// initialize the out object
	out := &pb.ControlResponse{
		Control: myMsg.(*pb.AvatarControl),
	}
	out.Control.CreatedAt = utils.GetTimeStamp()
	out.Control.CurrentStep = in.Control.CurrentStep
	out.Control.StepStatuses = s.avatar.GetStepStatusesObj()
	// check if this is an restart from UI
	currentStep := in.Control.CurrentStep
	if currentStep >= 0 && currentStep < int32(len(in.Control.StepStatuses)) {
		currentStepObj := in.Control.StepStatuses[currentStep]
		if currentStepObj.Status == pb.AvatarStatus_STATUS_SUCCESS {
			logrus.Infof("Continue from step %d", currentStep+1)
			logrus.Infof("Response from UI %s", currentStepObj.Response)
			// check if response is set

			if currentStepObj.Response != "" {
				// mutate the context, so LLM knows what the user sent
				if currentStepObj.ResponseType == pb.StepStatus_IMAGEDATA {

					// if the response type is IMAGEDATA
					// save the file to uploads folder
					// and push the filename into llm memory context

					newFileName := fmt.Sprintf("%s.%s", uuid.New().String(), "png")
					err := utils.SaveImageDataToFile(newFileName, currentStepObj.Response)
					if err != nil {
						return err
					}
					s.avatar.Memory.Messages = append(s.avatar.Memory.Messages, openai.UserMessage("image is stored in filename ```"+newFileName+"```"))
				} else {
					s.avatar.Memory.Messages = append(s.avatar.Memory.Messages, openai.UserMessage(currentStepObj.Response))
				}

			}
			in.Control.CurrentStep++
		}
	}

	// intimate client on starting the Avatar
	out.Control.Status = pb.AvatarStatus_STATUS_ACTIVE
	if err := stream.Send(out); err != nil {
		logrus.Errorf("Connect: Failed to send message: %v", err)
		return err
	}

	lastResponse := &pbConscious.ContextMesssage{
		Uuid: in.Control.Meta.Uuid,
	}

	for i := int(in.Control.CurrentStep); i < len(s.avatar.Avatar.Steps); i++ {
		step := s.avatar.Avatar.Steps[i]
		// for i, step := range s.avatar.Avatar.Steps {

		// intimate client on the step completed and its status
		out.Control.CurrentStep = int32(i)
		out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_ACTIVE
		out.Control.StepStatuses[i].CreatedAt = utils.GetTimeStamp()
		if err := stream.Send(out); err != nil {
			logrus.Errorf("Connect: Failed to send message: %v", err)
			return err
		}

		// run the actual step using llm
		outContextMsg, err := s.avatar.RunStep(i, step)

		out.Control.StepStatuses[i].CompletedAt = utils.GetTimeStamp()
		if err != nil {
			out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_FAILED
			out.Control.Status = pb.AvatarStatus_STATUS_FATAL
		} else {
			out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_SUCCESS
		}

		if outContextMsg.Action == "get_ui_launch_tool" {
			// we need to launch an UI in avatar
			out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_UI_WAITING
			out.Control.Status = pb.AvatarStatus_STATUS_UI_WAITING
			if err := stream.Send(out); err != nil {
				logrus.Errorf("Connect: Failed to send message: %v", err)
				return err
			}
			return nil
		}

		if outContextMsg.ResponseImage != "" {
			lastResponse.ResponseImage = outContextMsg.ResponseImage
			out.Control.StepStatuses[i].Response = outContextMsg.ResponseImage
		} else {
			lastResponse.Response = outContextMsg.Response
			out.Control.StepStatuses[i].Response = outContextMsg.Response
		}

		if err := stream.Send(out); err != nil {
			logrus.Errorf("Connect: Failed to send message: %v", err)
			return err
		}

	}
	logrus.Infof("Avatar run completed")
	out.Control.CompletedAt = utils.GetTimeStamp()
	out.Control.Status = pb.AvatarStatus_STATUS_SUCCESS
	out.Control.Response = lastResponse
	if err := stream.Send(out); err != nil {
		logrus.Errorf("Connect: Failed to send message: %v", err)
		return err
	}
	return nil
}
