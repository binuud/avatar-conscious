// consciousServer implements the ConsciousServer for the service from conscious proto file

// This server, listens on In endpoint, for user prompts, then makes a connection with backend AI,
// gets the response and returns same.
// When it gets a response from ai model, it publishes same to MQTT broker, so tts and other services, can listen to response events
// and respond for same
// States are saved, so each time the server is started, the old chat is maintained

package consciousServer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/utils"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"
	"github.com/google/uuid"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/sirupsen/logrus"
)

type grpcConsciousServer struct {
	pb.UnimplementedConsciousServer
	// memory       *memory.Memory
	openAiClient *openai.Client
}

const (
	DEFAULT_MEMORY_NAME = "default"
)

func NewGrpcConsciousServer() *grpcConsciousServer {

	// memory := memory.NewMemory(DEFAULT_MEMORY_NAME)

	// // override model from load with settings value
	// memory.Meta.Model = settings.Global.GetValue("LLM_MODEL")
	// memory.ReadSystemPromptFromSettings()

	// logrus.Infof("Starting Conscious with %s model", memory.Meta.Model)
	logrus.Infof("Starting Conscious")

	client := openai.NewClient(
		option.WithBaseURL(llm.GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	grpcServer := &grpcConsciousServer{
		// memory:       memory,
		openAiClient: &client,
	}

	return grpcServer
}

func (s *grpcConsciousServer) GetModels(ctx context.Context, in *pb.ModelRequest) (*pb.ModelResponse, error) {

	models, err := llm.MakeLLMGetModels()
	if err != nil {
		logrus.Errorf("Cannot get models %v", err)
		return nil, err
	}
	return &pb.ModelResponse{
		Models: models,
	}, nil

}

func (s *grpcConsciousServer) In(ctx context.Context, in *pb.ContextMesssage) (*pb.ContextMesssage, error) {

	in.CreatedAt = utils.GetTimeStamp()
	logrus.Infof("Generate Received -> Request: %s", in.Prompt)

	if in.Prompt == "" {
		return nil, fmt.Errorf("error prompt emtpy")
	}

	if in.Model == "" {
		return nil, fmt.Errorf("error - no model provided")
	}

	// do not add to memory for slash commands
	if strings.HasPrefix(in.Prompt, "/") {
		out := handleSlashCommand(in)
		return out, nil
	}

	r := NewRouter(in)
	out, err := r.RouteRequest(in)
	if err != nil {
		logrus.Errorf("LLM Error %v", err)
		return nil, err
	}

	logrus.Infof("In Received -> Response: %s", out)
	out.CompletedAt = utils.GetTimeStamp()
	return out, nil

}

func (s *grpcConsciousServer) InStream(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	r := NewRouter(in)
	return r.RouteStreamingPromptIntent(in, stream)

}

// Generate()
// uses tool handler
// tool handler to checnk if intent of user is to create image or describe an image
// was slowing the system, so using a slash command handler, user can specify his intent by giving
// loads temperature, seed settings from memory
// /image create /image edit or /image describe
func (s *grpcConsciousServer) Generate(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	logrus.Info("")
	logrus.Info("")
	logrus.Infof("Generate Received          -----------> Request: %s", in.Prompt)

	if in.Prompt == "" {
		return fmt.Errorf("error prompt emtpy")
	}

	// do not add to memory for slash commands
	if strings.HasPrefix(in.Prompt, "/") {
		return s.handleSlashCommandStream(in, stream)
	}

	r := NewRouter(in)
	err := r.RouteStreamingPromptIntent(in, stream)
	r.memory.DumpMemory(nil)
	return err
}

func (s *grpcConsciousServer) FileUpload(ctx context.Context, in *pb.UploadRequest) (*pb.UploadResponse, error) {

	logrus.Infof("Recieved file for upload %s", in.FileName)
	_, err := utils.IsFilenameSafe(in.FileName)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid file extension: "+err.Error())
	}

	if in.GetContentType() == "" {
		// Default to application/octet-stream if not provided
		in.ContentType = "application/octet-stream"
	}

	// get the extension of the file and add to the UUID
	// so new file is save as [uuid].png etc.
	ext := filepath.Ext(in.FileName)
	ext = strings.ToLower(ext)
	// Remove the dot from extension for the new filename
	ext = strings.TrimPrefix(ext, ".")
	newFileName := fmt.Sprintf("%s.%s", uuid.New().String(), ext)

	err = utils.SaveImageDataToFile(newFileName, in.FileData)
	if err != nil {
		return nil, err
	}

	return &pb.UploadResponse{
		UploadId:     newFileName,
		Status:       "success",
		ErrorMessage: "",
	}, nil
}

func (s *grpcConsciousServer) SensorEvents(stream grpc.BidiStreamingServer[pb.ContextMesssage, pb.ContextMesssage]) error {

	logrus.Debug("Connect: Client connected to Conscious SensorEvents bidirectional stream")
	return s.handleSensorEvents(stream)

}

// // SetModel
// // SetModel in the current chat session
// func (s *grpcConsciousServer) SetModel(ctx context.Context, in *pb.SetModelRequest) (*pb.SetModelResponse, error) {

// 	logrus.Info("Recieved SetModel")
// 	s.memory.SetModel(in.Model)
// 	s.memory.Save()

// 	return &pb.SetModelResponse{Model: s.memory.Meta.Model}, nil
// }

// ReloadTools
// Deletes the memory from storage and creates new
// Handle mutex ##TBD Binu
func (s *grpcConsciousServer) ReloadTools(ctx context.Context, in *pb.ReloadToolsRequest) (*pb.ReloadToolsResponse, error) {

	logrus.Info("Recieved ReloadTools")
	if in.MemoryName == "" {
		return nil, status.Error(codes.InvalidArgument, "emtpy memory name ")
	}
	memory := memory.NewMemory(in.MemoryName)
	memory.InitTools()
	memory.Save()
	memory.DumpTools()
	return &pb.ReloadToolsResponse{
		Success: true,
	}, nil
}

func (s *grpcConsciousServer) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthResponse, error) {

	health := &pb.HealthResponse{
		Uuid:          uuid.New().String(),
		Status:        pb.HealthResponse_Active,
		StatusMessage: "Need to integrate with api TODO",
		CreatedAt:     time.Now().Unix(),
	}

	return health, nil
}
