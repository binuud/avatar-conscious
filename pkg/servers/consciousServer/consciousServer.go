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
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/binuud/avatar-conscious/pkg/utils"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"
	"github.com/google/uuid"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"

	"github.com/sirupsen/logrus"
)

type grpcConsciousServer struct {
	pb.UnimplementedConsciousServer
	memory       *memory.Memory
	openAiClient *openai.Client
	commonTools  []openai.ChatCompletionToolUnionParam
}

const (
	DEFAULT_MEMORY_NAME = "default"
)

func NewGrpcConsciousServer() *grpcConsciousServer {

	memory := memory.NewMemory(DEFAULT_MEMORY_NAME)

	logrus.Infof("Starting Conscious with %s model", memory.Model)

	client := openai.NewClient(
		option.WithBaseURL(llm.GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	grpcServer := &grpcConsciousServer{
		memory:       memory,
		openAiClient: &client,
	}
	grpcServer.commonTools = tools.GlobalCommonTools

	datalayer.SaveJSON(datalayer.GetFileAbsolutePath("configs", "config.json"), grpcServer.commonTools)

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

	logrus.Infof("In Received -> Request: %s, %s", in.Type, in.Prompt)

	out, err := s.RouteRequest(in)
	if err != nil {
		logrus.Errorf("LLM Error %v", err)
		return nil, err
	}

	logrus.Infof("In Received -> Response: %s", out)
	return out, nil

}

func (s *grpcConsciousServer) InStream(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	s.memory.AppendInMessage(in)
	// thinking := ""
	completeResponse := ""

	s.memory.AppendInMessage(in)

	llmStream := llm.MakeLLMChatStreamingChatRequest(s.memory.Messages, s.memory.Tools, s.memory.Model, s.memory.Temperature)

	defer func() {
		_ = llmStream.Close()
	}()

	logrus.Infof("InStream Received -> Request: %s, %s", in.Type, in.Prompt)

	// Iterate through the stream chunks as they arrive
	for llmStream.Next() {
		evt := llmStream.Current()
		if len(evt.Choices) > 0 {
			// Print each chunk to stdout in real-time
			fmt.Print(evt.Choices[0].Delta.Content)
			out := &pb.ContextMesssage{
				Uuid:     in.Uuid,
				Response: evt.Choices[0].Delta.Content,
			}
			completeResponse = completeResponse + evt.Choices[0].Delta.Content
			// 	// send the out message
			if err := stream.Send(out); err != nil {
				return err
			}
		}
	}
	logrus.Info("Stream closed")

	// if streaming is done
	// send last out message with completeAt timestamp
	out := &pb.ContextMesssage{
		Uuid:        in.Uuid,
		Model:       in.Model,
		Prompt:      in.Prompt,
		CreatedAt:   in.CreatedAt,
		CompletedAt: utils.GetTimeStamp(),
		Response:    completeResponse,
	}

	s.memory.AppendOutMessage(out)

	// send the out message
	if err := stream.Send(out); err != nil {
		return err
	}

	return nil

}

func (s *grpcConsciousServer) InStreamNewApi(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	logrus.Infof("InStream Received -> Request: %s, %s", in.Type, in.Prompt)

	// thinking := ""
	completeResponse := ""

	s.memory.AppendInMessage(in)

	params := responses.ResponseNewParams{
		Model: in.Model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentParamOfInputText(in.Prompt),
					},
					responses.EasyInputMessageRoleUser,
				),
			},
		},
	}

	llmStream := s.openAiClient.Responses.NewStreaming(context.Background(), params)
	defer func() {
		_ = llmStream.Close()
	}()

	logrus.Info("\n\n\n Starting streaming")

	// Iterate over incoming chunks as they arrive from the API
	for llmStream.Next() {

		if err := stream.Context().Err(); err != nil {
			logrus.Errorf("Client disconnected prematurely. %v", err)
			return err
		}

		data := llmStream.Current()

		out := &pb.ContextMesssage{
			Uuid:     in.Uuid,
			Response: data.Delta,
		}

		if data.JSON.Text.Valid() {

			// when reasoning, thinking is enabled, we might see this twice.
			// once for reasoning, and once for actualy response

			logrus.Infof("Finished Content %s \n\n", data.Type)

			if data.Type == "response.output_text.done" {
				// streaming complete
				// handle closure later
				completeResponse = data.Text
				logrus.Infof("Complete response %s \n\n", completeResponse)
			}
		}

		// send the out message
		if err := stream.Send(out); err != nil {
			return err
		}

	}

	// if streaming is done
	// send last out message with completeAt timestamp
	out := &pb.ContextMesssage{
		Uuid:        in.Uuid,
		Model:       in.Model,
		Prompt:      in.Prompt,
		CreatedAt:   in.CreatedAt,
		CompletedAt: utils.GetTimeStamp(),
		Response:    completeResponse,
	}

	s.memory.AppendOutMessage(out)

	// send the out message
	if err := stream.Send(out); err != nil {
		return err
	}

	return nil

}

// Generate()
// uses new tool handler
// Calls toools first with no memory, to decide if image generation or description or tts has to be called
// Then calls the appropriate function
func (s *grpcConsciousServer) Generate(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	in.CreatedAt = utils.GetTimeStamp()
	logrus.Infof("Generate Received -> Request: %s", in.Prompt)

	// do not add to memory for slash commands
	if strings.HasPrefix(in.Prompt, "/") {
		return handleSlashCommand(in, stream)
	}

	// #TODO Binu hot fix shitty hack
	// should be removed and cleaner implementation done
	if in.Data != "" {
		// add the attached file name to the prompt
		in.Prompt = in.Prompt + " . And the attached or uploaded file is " + in.Data
	}

	// for tool calling, do not send the history
	message := []openai.ChatCompletionMessageParamUnion{openai.UserMessage(in.Prompt)}
	completions, err := llm.MakeLLMSingleToolRequest(message,
		s.memory.Model,
		s.commonTools,
		s.memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return err
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		s.memory.AppendInMessage(in)
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return err
		}
		// s.memory.Messages = append(s.memory.Messages, openai.AssistantMessage(toolMessage))

		// for metrics, populate the time stamps
		out.CreatedAt = in.CreatedAt
		out.CompletedAt = utils.GetTimeStamp()

		logrus.Infof("Generate Out Msg : %s", out)
		s.memory.AppendOutMessage(out)
		if err := stream.Send(out); err != nil {
			return err
		}
		return nil
	} else {

		logrus.Infof("Generate Tool function NOT selected\n")
		logrus.Infof("Generate Out Msg : %s", msg.Content)
		if in.Data != "" {
			// if an image is in the attachment, do a describe image
			in.Type = pb.ContextMesssage_LLM_IMAGE_DESCRIBE
			in.Model = settings.Global.GetValue("VISION_MODEL")
		}
	}

	// this is a normal chat streaming call
	// user has not asked for any tools
	return s.InStream(in, stream)

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

// ClearMemory
// Deletes the memory from storage and creates new
// Handle mutex ##TBD Binu
func (s *grpcConsciousServer) ClearMemory(ctx context.Context, in *pb.ClearMemoryRequest) (*pb.ClearMemoryResponse, error) {

	memory.DeleteMemory(DEFAULT_MEMORY_NAME)
	s.memory = memory.NewMemory(DEFAULT_MEMORY_NAME)

	return &pb.ClearMemoryResponse{}, nil
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

func (s *grpcConsciousServer) MakeLLMChatRequest(prompt string) (*openai.ChatCompletion, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(prompt))
	completions, err := llm.MakeLLMChatRequest(s.memory.Messages, s.memory.Tools, s.memory.Model, s.memory.Temperature)
	if err != nil {
		return nil, err
	}

	return completions, nil

}
