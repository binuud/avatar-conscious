package avatarServer

import (
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

func (s *avatarGrpcServer) handleAvatarBiDirectionalServer(stream grpc.BidiStreamingServer[pb.ControlRequest, pb.ControlResponse]) error {
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
		logrus.Infof("Received action from client %v, Step No: %d ", in.Control.Action, (in.Control.CurrentStep + 1))

		switch in.Control.Action {
		// initialize the avatar session
		case pb.AvatarAction_ACTION_INIT:
			{
				err := s.handleAvatarInit(in, stream)
				if err != nil {
					return err
				}
			}
		// start the avatar and run to completion
		// break if ui intervention is needed
		case pb.AvatarAction_ACTION_BEGIN:
			{
				err := s.handleAvatarRun(in, stream)
				if err != nil {
					return err
				}
			}

		// run avatar step by step
		// in debug mode
		case pb.AvatarAction_ACTION_DEBUG:
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
		case pb.AvatarAction_ACTION_STOP:
			{
				logrus.Infof("Avatar %s, stopped on user request", s.avatar.Meta.Name)
				out.Control.Status = pb.AvatarStatus_STATUS_ERROR
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
	s.avatar = NewAvatar(in.Control.Name)

	// initialize memory
	s.avatar.Start()
	myMsg := proto.Clone(in.Control)
	// initialize the out object
	out := &pb.ControlResponse{
		Control: myMsg.(*pb.AvatarControl),
	}
	out.Control.CurrentStep = 0
	if s.avatar == nil {
		out.Control.Status = pb.AvatarStatus_STATUS_ERROR
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
		userResponse := currentStepObj.Response
		prompt := ""
		if currentStepObj.Status == pb.AvatarStatus_STATUS_COMPLETED {
			logrus.Infof("Continue from step %d", currentStep+1)
			logrus.Infof("Response from UI %s", userResponse[:min(len(userResponse), 100)])
			// check if response is set

			if userResponse != "" {
				// mutate the context, so LLM knows what the user sent
				if currentStepObj.ResponseType == pb.StepStatus_IMAGEDATA {

					// if the response type is IMAGEDATA
					// save the file to uploads folder
					// and push the filename into llm memory context

					newFileName := fmt.Sprintf("%s.%s", uuid.New().String(), "png")
					err := utils.SaveImageDataToFile(newFileName, userResponse)
					if err != nil {
						return err
					}
					prompt = fmt.Sprintf("Step No: %d, image is stored in filename ```%s```", (currentStep + 1), newFileName)

				} else {
					prompt = fmt.Sprintf("Step No: %d, %s", (currentStep + 1), currentStepObj.Response)
				}

				s.avatar.Memory.AppendMemory(openai.UserMessage(prompt))
				s.avatar.Memory.AppendMemory(openai.AssistantMessage("Ok."))

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
		Uuid: in.Control.AvatarUuid,
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
		if err != nil {
			out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_ERROR
			out.Control.StepStatuses[i].Result = pb.AvatarResult_RESULT_ERROR
			out.Control.Status = pb.AvatarStatus_STATUS_ERROR
			out.Control.Result = pb.AvatarResult_RESULT_ERROR
		} else {
			out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_COMPLETED
		}
		out.Control.StepStatuses[i].CompletedAt = utils.GetTimeStamp()

		if outContextMsg.Action == "SHOW-UI" {
			// we need to launch an UI in avatar
			out.Control.StepStatuses[i].Status = pb.AvatarStatus_STATUS_WAITING
			out.Control.StepStatuses[i].Response = outContextMsg.Data
			out.Control.Status = pb.AvatarStatus_STATUS_WAITING

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
	out.Control.Action = pb.AvatarAction_ACTION_COMPLETE
	out.Control.Status = pb.AvatarStatus_STATUS_COMPLETED
	out.Control.Result = pb.AvatarResult_RESULT_SUCCESS
	out.Control.Response = lastResponse
	if err := stream.Send(out); err != nil {
		logrus.Errorf("Connect: Failed to send message: %v", err)
		return err
	}
	return nil
}
