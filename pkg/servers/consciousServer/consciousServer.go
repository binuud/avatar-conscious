package consciousServer

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	protoV1 "github.com/binuud/conscious/gen/go/v1/brain"
	"github.com/openai/openai-go"
	"github.com/sirupsen/logrus"
)

type grpcConsciousServer struct {
	protoV1.UnimplementedConsciousServer
	model        string
	systemPrompt string
	temperature  float64
	messages     []openai.ChatCompletionMessageParamUnion
}

func NewGrpcConsciousServer() *grpcConsciousServer {

	systemPrompt := "You are a very helpful assistant. Very polite, friendly with children. Do not use curse words. Do not use bad words. Reply briefly, in less than 20 words as far as possible. Do not return smiley, icons or images."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	return &grpcConsciousServer{
		model:        "gemma3",
		systemPrompt: systemPrompt,
		temperature:  0.8,
		messages:     messages,
	}

}

func (s *grpcConsciousServer) In(ctx context.Context, in *protoV1.InRequest) (*protoV1.InResponse, error) {

	response, err := s.MakeLLMChatRequest(in.Prompt)
	if err != nil {
		logrus.Errorf("LLM Error %v", err)
		return nil, err
	}

	return &protoV1.InResponse{
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

func (s *grpcConsciousServer) MakeLLMChatRequest(prompt string) (string, error) {

	s.messages = append(s.messages, openai.UserMessage(prompt))
	completions, err := MakeLLMChatRequest(prompt, s.messages, s.model, s.temperature)
	if err != nil {
		return "", err
	}

	// add response to chat
	aiMsg := completions.Choices[0].Message.Content
	fmt.Println("AI:", aiMsg)
	s.messages = append(s.messages, openai.AssistantMessage(aiMsg))

	return aiMsg, nil

}
