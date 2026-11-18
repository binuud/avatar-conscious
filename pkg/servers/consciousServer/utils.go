package consciousServer

import (
	"bytes"
	"strings"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/slashCommands"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

func (s *grpcConsciousServer) MakeLLMChatRequest(prompt string) (*openai.ChatCompletion, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(prompt))
	completions, err := llm.MakeLLMChatRequest(s.memory.Messages, s.memory.Tools, s.memory.Meta.Model, s.memory.Meta.Temperature)
	if err != nil {
		return nil, err
	}

	return completions, nil

}

// handleTool
// returns a bool, and error. The boolean is to indicate whether the tool was handled
func (s *grpcConsciousServer) handleTool(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) (bool, *pb.ContextMesssage, error) {

	// adding user prompt, so that the tooling call can decide what tool to use
	s.memory.AppendInMessage(in)
	completions, err := llm.MakeLLMSingleToolRequest(s.memory.Messages,
		s.memory.Meta.Model,
		s.memory.Tools,
		s.memory.Meta.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return false, nil, err
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {

		out, err := s.memory.ToolRegistry.HandleToolCalls(msg.ToolCalls, in.Prompt)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return false, nil, err
		}
		toolMessage := out.Response
		if out.ResponseImage != "" {
			toolMessage = out.ResponseImage
		}
		s.memory.AppendMessage(openai.ToolMessage(toolMessage, msg.ToolCalls[0].ID))
		// for metrics, populate the time stamps
		out.CreatedAt = in.CreatedAt
		out.CompletedAt = utils.GetTimeStamp()

		logrus.Debugf("Generate Out Msg : %s", out)
		// s.memory.AppendOutMessage(out)
		if err := stream.Send(out); err != nil {
			return true, out, err
		}
		return true, out, nil
	} else {

		// logrus.Debugf("Generate Tool function NOT selected\n")
		// logrus.Debugf("Generate Out Msg : %s", msg.Content)
		// if in.Data != "" {
		// 	// if an image is in the attachment, do a describe image
		// 	in.Type = pb.ContextMesssage_LLM_IMAGE_DESCRIBE
		// 	in.Model = settings.Global.GetValue("VISION_MODEL")
		// }
	}
	return false, nil, nil
}

// handleSlashCommandStream, first checks if the slash command is in the list
// image edit, image create, image describe, if so it handles them
// if the slash command does not fit in the above list, it is executed by the slash command handler
// note, tool routing was taking from 1 to 10 seconds with big context sizes
// so not use tool routing to check if the prompt is for image creation, edit or delete
// if no suitable prompt intent is found, it tries to handle the tool calling
// if no tools are found, then it tries to hanndle the prompt using llm
func (s *grpcConsciousServer) handleSlashCommandStream(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {
	// if prompt starts with a "/" treat is as a command, and not a prompt

	logrus.Infof("Handling slash command %s", in.Prompt)

	out := &pb.ContextMesssage{}

	var err error

	if strings.HasPrefix(in.Prompt, "/image describe") {
		in.Model = settings.Global.GetValue("VISION_MODEL")
		out, err = s.RouteLLMImageDescribeRequest(in)
	} else if strings.HasPrefix(in.Prompt, "/image create") {

		// for image create, use the global image model
		// if user needs to select model, then the user has to selec the IMAGE CREATE type request
		in.Model = settings.Global.GetValue("IMAGE_MODEL")
		out, err = s.RouteOpenAiImageGenerationRequest(in)

	} else if strings.HasPrefix(in.Prompt, "/image edit") {

		// for image edit, use the global image model
		// if user needs to select model, then the user has to selec the IMAGE EDIT type request
		in.Model = settings.Global.GetValue("IMAGE_MODEL")
		out, err = s.RouteOpenAiImageEditRequest(in)

	} else {
		out = handleSlashCommand(in)
	}

	if err != nil {
		logrus.Errorf("err in image edit/describe %v", err)
		out.Response = err.Error()
		if err = stream.Send(out); err != nil {
			return err
		}
	}

	if err = stream.Send(out); err != nil {
		return err
	}

	return nil
}

func handleSlashCommand(in *pb.ContextMesssage) *pb.ContextMesssage {
	var buf bytes.Buffer
	slashCommands.SetStream(&buf)
	slashCommands.HandleSlashCommand(in.Prompt[1:])
	return &pb.ContextMesssage{
		Response: buf.String(),
	}

}
