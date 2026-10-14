// consciousServer implements the ConsciousServer for the service from conscious proto file

// This server, listens on In endpoint, for user prompts, then makes a connection with backend AI,
// gets the response and returns same.
// When it gets a response from ai model, it publishes same to MQTT broker, so tts and other services, can listen to response events
// and respond for same
// States are saved, so each time the server is started, the old chat is maintained

package consciousServer

import (
	"context"
	"os"
	"time"

	pb "github.com/binuud/conscious/gen/go/v1/conscious"
	grpc "google.golang.org/grpc"

	// mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"
	"github.com/google/uuid"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"

	"github.com/sirupsen/logrus"
)

type memory struct {
	Name         string
	Model        string
	SystemPrompt string
	Temperature  float64
	Messages     []openai.ChatCompletionMessageParamUnion
	Tools        []openai.ChatCompletionToolUnionParam
}

type grpcConsciousServer struct {
	pb.UnimplementedConsciousServer
	memory *memory
}

func NewGrpcConsciousServer() *grpcConsciousServer {

	systemPrompt := "You are a helpful weather assistant. You must use the tools provided to answer questions. Very polite, friendly with children. Do not use curse words. Do not use bad words. Reply briefly, in less than 20 words as far as possible. Do not use any smilies, emojis, or symbols in your responses from now on."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	LLM_MODEL := os.Getenv("LLM_MODEL")
	if LLM_MODEL == "" {
		LLM_MODEL = "gemma3"
	}
	tools := []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_weather",
			Description: openai.String("Get the current weather for a given city or place"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]string{
						"type":        "string",
						"description": "The city or place, e.g. Tokyo",
					},
				},
				"required": []string{"location"},
			},
		}),
	}

	memory := &memory{
		Name:         "Default",
		Model:        LLM_MODEL,
		SystemPrompt: systemPrompt,
		Temperature:  0.8,
		Messages:     messages,
	}

	err := memory.Load("default")

	// do not save tools because of issue with saving hash map
	memory.Tools = tools
	memory.DumpMemory()
	if err != nil {
		logrus.Infof("No Default memory found")
	} else {
		logrus.Infof("Default memory found and loaded")
	}
	logrus.Infof("Starting Conscious with %s model", LLM_MODEL)
	return &grpcConsciousServer{memory: memory}

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

	err = s.memory.Save("default")
	if err != nil {
		logrus.Panicf("Cannot save conscious memory %v", err)
	}
	logrus.Infof("In Received -> Response: %s", out)
	return out, nil

}

// Removed the 'ctx context.Context' parameter and fixed the stream interface type
func (s *grpcConsciousServer) InStream(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	logrus.Infof("InStream Received -> Request: %s, %s", in.Type, in.Prompt)
	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	// logrus.Infof("Send data to llm server %s", params.User.Value)

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

	llmStream := client.Responses.NewStreaming(context.Background(), params)
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
