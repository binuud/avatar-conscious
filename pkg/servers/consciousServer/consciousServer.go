package consciousServer

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"

	protoV1 "github.com/binuud/conscious/gen/go/v1/brain"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
)

type grpcConsciousServer struct {
	protoV1.UnimplementedConsciousServer
}

func NewGrpcConsciousServer() *grpcConsciousServer {

	return &grpcConsciousServer{}

}

func (s *grpcConsciousServer) In(ctx context.Context, in *protoV1.Request) (*protoV1.Response, error) {

	response, err := MakeLLMRequest(in.Prompt)
	if err != nil {
		logrus.Errorf("LLM Error %v", err)
		return nil, err
	}

	return &protoV1.Response{
		Uuid:     in.Uuid,
		Prompt:   in.Prompt,
		Response: response,
	}, nil

}

func (s *grpcConsciousServer) Health(ctx context.Context, in *protoV1.HealthRequest) (*protoV1.HealthResponse, error) {

	health := &protoV1.HealthResponse{}
	health.Uuid = uuid.New().String()
	// check if the rag vector is populated

	health.Status = protoV1.HealthResponse_Active
	health.StatusMessage = "Need to integrate with api TODO"

	return health, nil
}

func MakeLLMRequest(prompt string) (string, error) {
	llmURL := os.Getenv("LLM_ENDPOINT") + "/engines/llama.cpp/v1/"
	model := "gemma3"

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a useful AI agent expert with TV series."),
		openai.UserMessage(prompt),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(0.8),
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		return "", err
	}

	var response = completion.Choices[0].Message.Content

	fmt.Println(response)
	return response, nil
}
