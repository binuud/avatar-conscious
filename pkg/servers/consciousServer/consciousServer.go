// consciousServer implements the ConsciousServer for the service from conscious proto file

// This server, listens on In endpoint, for user prompts, then makes a connection with backend AI,
// gets the response and returns same.
// When it gets a response from ai model, it publishes same to MQTT broker, so tts and other services, can listen to response events
// and respond for same
// States are saved, so each time the server is started, the old chat is maintained

package consciousServer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/binuud/conscious/pkg/utils"
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
	memory       *Memory
	openAiClient *openai.Client
	commonTools  []openai.ChatCompletionToolUnionParam
}

func NewGrpcConsciousServer() *grpcConsciousServer {

	systemPrompt := "You are a helpful assistant."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	LLM_MODEL := os.Getenv("LLM_MODEL")
	if LLM_MODEL == "" {
		LLM_MODEL = "gemma3"
	}

	memory := &Memory{
		Name:         "Default",
		Model:        LLM_MODEL,
		SystemPrompt: systemPrompt,
		Temperature:  0.8,
		Messages:     messages,
	}

	err := memory.Load()

	// do not save tools because of issue with saving hash map
	memory.DumpMemory()
	if err != nil {
		logrus.Infof("No Default memory found")
	} else {
		logrus.Infof("Default memory found and loaded")
	}
	logrus.Infof("Starting Conscious with %s model", LLM_MODEL)

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	grpcServer := &grpcConsciousServer{
		memory:       memory,
		openAiClient: &client,
	}
	grpcServer.commonTools = grpcServer.addCommonToolsRoute()

	return grpcServer
}

func (s *grpcConsciousServer) addCommonToolsRoute() []openai.ChatCompletionToolUnionParam {
	return []openai.ChatCompletionToolUnionParam{
		// tool case for image generation
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_image",
			Description: openai.String("Retreive or create an image"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"text": map[string]string{
						"type":        "string",
						"description": "A description of the image",
					},
				},
				"required": []string{"text"},
			},
		}),
		// tool case for image description
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_image_description",
			Description: openai.String("Describe a given image or photo."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"text": map[string]string{
						"type":        "string",
						"description": "A question about the image",
					},
				},
				"required": []string{"text"},
			},
		}),
	}
}

func (s *grpcConsciousServer) GetModels(ctx context.Context, in *pb.ModelRequest) (*pb.ModelResponse, error) {

	models, err := MakeLLMGetModels()
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

	logrus.Infof("InStream Received -> Request: %s, %s", in.Type, in.Prompt)

	s.memory.AppendMemory(openai.UserMessage(in.Prompt))

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
		// fmt.Printf("Type %s, data %s \n", data.Type, data.Delta)
		resp := &pb.ContextMesssage{
			Response: data.Delta,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		if data.JSON.Text.Valid() {

			// when reasoning, thinking is enabled, we might see this twice.
			// once for reasoning, and once for actualy response
			logrus.Infof("Finished Content %s \n\n", data.Type)
			// completeText = data.Text
		}
	}
	return nil
}

// Generate,
// Calls toools first with no memory, to decide if image generation or description or tts has to be called
// Then calls the appropriate function
func (s *grpcConsciousServer) Generate(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	logrus.Infof("Generate Received -> Request: %s", in.Prompt)

	// make a tool request
	// LLM will decide if the user wants to create image, create tts,
	// describe image or chat stream
	logrus.Infof("Tools %v\n", s.commonTools)
	tool, err := GetLLMToolName(in.Prompt, "", in.Model, s.commonTools, s.memory.Temperature)
	if err != nil {
		// t.Errorf("client: error llm chat request: %v \n", err)
		return err
	}

	// if llm model returns a tool name
	if tool != nil {
		logrus.Infof("Tool function selected : %s with arguments: %s\n", tool.Function.Name, tool.Function.Arguments)

		switch tool.Function.Name {
		case "get_image":
			in.Type = pb.ContextMesssage_IMAGE_GENERATION
			in.Model = utils.GetEnv("IMAGE_MODEL", "Runpod/FLUX.2-klein-4B-mflux-4bit")
		case "get_image_description":
			in.Type = pb.ContextMesssage_LLM_IMAGE_DESCRIBE
			in.Model = utils.GetEnv("VISION_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit")
		default:
			logrus.Errorf("Invalid tool name selected by llm model %s", tool.Function.Name)
			return fmt.Errorf("Invalid tool selected by llm model %s", tool.Function.Name)
		}

		// if model returns refined prompt,
		// hot swap it for user prompt
		var args struct {
			Text string `json:"text"`
		}
		json.Unmarshal([]byte(tool.Function.Arguments), &args)
		newPrompt := args.Text

		if newPrompt != "" {
			in.Prompt = newPrompt
		}

		// call the router function
		// and modify the prompt

		out, err := s.RouteRequest(in)
		if err != nil {
			logrus.Errorf("LLM Error %v", err)
			return err
		}

		err = s.memory.Save()
		if err != nil {
			logrus.Panicf("Cannot save conscious memory %v", err)
		}
		logrus.Infof("In Received -> Response: %s", out)

		// return the output response and exit stream
		if err := stream.Send(out); err != nil {
			return err
		}
		return nil
	} else {
		logrus.Infof("Tool function NOT selected\n")

		if in.Data != "" {
			// if an image is in the attachment, do a describe image
			in.Type = pb.ContextMesssage_LLM_IMAGE_DESCRIBE
			in.Model = utils.GetEnv("VISION_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit")
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

	rawBase64 := in.FileData

	// 1. Fix Corruption: Locate and strip the "data:image/png;base64," prefix if it exists
	if idx := strings.Index(rawBase64, ","); idx != -1 {
		rawBase64 = rawBase64[idx+1:]
	}

	// 2. Decode the clean, sanitized base64 payload string
	decodedBytes, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot decode file: "+err.Error())
	}

	if len(decodedBytes) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file empty")
	}

	ext := filepath.Ext(in.FileName)
	ext = strings.ToLower(ext)
	// Remove the dot from extension for the new filename
	ext = strings.TrimPrefix(ext, ".")
	newFileName := fmt.Sprintf("%s.%s", uuid.New().String(), ext)

	destinationPath := utils.GetEnv("FILE_STORE_PATH", ".") + "/uploads/" + newFileName
	logrus.Infof("Uplaoded file to be saveed under -> %s", destinationPath)

	err = os.WriteFile(destinationPath, decodedBytes, 0644)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot write file to server: "+err.Error())
	}

	return &pb.UploadResponse{
		UploadId:     newFileName,
		Status:       "success",
		ErrorMessage: "",
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

func (s *grpcConsciousServer) MakeLLMChatRequest(prompt string) (*openai.ChatCompletion, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(prompt))
	completions, err := MakeLLMChatRequest(prompt, s.memory.Messages, s.memory.Tools, s.memory.Model, s.memory.Temperature)
	if err != nil {
		return nil, err
	}

	return completions, nil

}
