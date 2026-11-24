package consciousServer

import (
	"encoding/json"
	"errors"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/servers/mcpManagerServer"
	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

const (
	OBJECT_IDENTIFICATION_MCP_SERVER = "yolo-mcp"
	OBJECT_IDENTIFICATION_MCP_TOOL   = "detect_objects_from_buffer"
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
		// logrus.Infof("Received sensor event from client %s(%s) ", in.Name, in.Type)

		out = &pb.ContextMesssage{
			Name:     in.Name,
			Response: "Unhandled",
		}

		if in.Type == pb.ContextMesssage_OBJECT_DETECTION {
			// handle mcp yolo
			out, err = s.handleObjectIdentification(in)
		}

		if err := stream.Send(out); err != nil {
			logrus.Errorf("Connect: Failed to send message: %v", err)
			return err
		}
	}
}

func (s *grpcConsciousServer) handleObjectIdentification(in *pb.ContextMesssage) (*pb.ContextMesssage, error) {

	// logrus.Debugf("handleObjectIdentification sensor events %v", in.Prompt[:100])

	inObj := make(mcpHandler.Tool_MCP_Input, 0)
	inObj["confidence_threshold"] = "0.25"
	inObj["image_base64"] = in.Prompt
	params, err := json.Marshal(inObj)
	if err != nil {
		return nil, err
	}

	mcp := mcpManagerServer.NewMcpServers()
	outMsg, err := mcp.HandleTool("", OBJECT_IDENTIFICATION_MCP_SERVER, OBJECT_IDENTIFICATION_MCP_TOOL, string(params))
	if err != nil {
		logrus.Errorf("handleObjectIdentification result %v", err)
		return nil, err
	}

	logrus.Debugf("handleObjectIdentification result %v", outMsg)
	out := &pb.ContextMesssage{
		Type:     in.Type,
		Name:     in.Name,
		Response: outMsg,
	}

	return out, nil
}
