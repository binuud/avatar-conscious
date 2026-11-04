package consciousServer

import (
	"bytes"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/slashCommands"
	grpc "google.golang.org/grpc"
)

func handleSlashCommandStream(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {
	// if prompt starts with a "/" treat is as a command, and not a prompt
	out := handleSlashCommand(in)
	if err := stream.Send(out); err != nil {
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
