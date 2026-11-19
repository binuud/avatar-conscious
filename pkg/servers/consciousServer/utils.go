package consciousServer

import (
	"bytes"
	"strings"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/slashCommands"
	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

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

	r := NewRouter(in)

	if strings.HasPrefix(in.Prompt, "/image describe") {
		in.Model = settings.Global.GetValue("VISION_MODEL")
		out, err = r.RouteLLMImageDescribeRequest(in)
	} else if strings.HasPrefix(in.Prompt, "/image create") {

		// for image create, use the global image model
		// if user needs to select model, then the user has to selec the IMAGE CREATE type request
		in.Model = settings.Global.GetValue("IMAGE_MODEL")
		out, err = r.RouteOpenAiImageGenerationRequest(in)

	} else if strings.HasPrefix(in.Prompt, "/image edit") {

		// for image edit, use the global image model
		// if user needs to select model, then the user has to selec the IMAGE EDIT type request
		in.Model = settings.Global.GetValue("IMAGE_MODEL")
		out, err = r.RouteOpenAiImageEditRequest(in)

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
