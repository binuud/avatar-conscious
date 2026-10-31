package consciousServer

import (
	"bytes"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/slashCommands"
	grpc "google.golang.org/grpc"
)

func handleSlashCommand(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {
	// if prompt starts with a "/" treat is as a command, and not a prompt
	var buf bytes.Buffer
	// remove the first slash and handle the command
	// s.memory.AppendInMessage(in)
	slashCommands.SetStream(&buf)
	slashCommands.HandleSlashCommand(in.Prompt[1:])
	out := &pb.ContextMesssage{
		Response: buf.String(),
	}

	// s.memory.AppendOutMessage(out)
	if err := stream.Send(out); err != nil {
		return err
	}
	return nil
}
