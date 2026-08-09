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
	"os"

	protoV1 "github.com/binuud/conscious/gen/go/v1/brain"
	aiModuleProto "github.com/dronasys-com/brain-nodes/gen/go/v1/aiModule"
	mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/sirupsen/logrus"
)

type memory struct {
	Name         string
	Model        string
	SystemPrompt string
	Temperature  float64
	Messages     []openai.ChatCompletionMessageParamUnion
}

type grpcConsciousServer struct {
	protoV1.UnimplementedConsciousServer
	memory *memory
}

func NewGrpcConsciousServer() *grpcConsciousServer {

	systemPrompt := "You are a very helpful assistant. Very polite, friendly with children. Do not use curse words. Do not use bad words. Reply briefly, in less than 20 words as far as possible. Do not use any smilies, emojis, or symbols in your responses from now on."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	LLM_MODEL := os.Getenv("LLM_MODEL")
	if LLM_MODEL == "" {
		LLM_MODEL = "gemma3"
	}

	memory := &memory{
		Name:         "Default",
		Model:        LLM_MODEL,
		SystemPrompt: systemPrompt,
		Temperature:  0.8,
		Messages:     messages,
	}

	err := memory.Load("default")
	memory.DumpMemory()
	if err != nil {
		logrus.Infof("No Default memory found")
	} else {
		logrus.Infof("Default memory found and loaded")
	}
	logrus.Infof("Starting Conscious with %s model", LLM_MODEL)
	return &grpcConsciousServer{memory: memory}

}

func (s *grpcConsciousServer) In(ctx context.Context, in *protoV1.InRequest) (*protoV1.InResponse, error) {

	response, err := s.MakeLLMChatRequest(in.Prompt)
	if err != nil {
		logrus.Errorf("LLM Error %v", err)
		return nil, err
	}
	mqttServer.SendOnTTS(response)

	err = s.memory.Save("default")
	if err != nil {
		logrus.Panicf("Cannot save conscious memory %v", err)
	}

	return &protoV1.InResponse{
		Uuid:     in.Uuid,
		Prompt:   in.Prompt,
		Response: response,
	}, nil

}

func (s *grpcConsciousServer) Health(ctx context.Context, in *aiModuleProto.HealthRequest) (*aiModuleProto.HealthResponse, error) {

	health := &aiModuleProto.HealthResponse{}
	health.Uuid = uuid.New().String()
	// check if the rag vector is populated

	health.Status = aiModuleProto.HealthResponse_Active
	health.StatusMessage = "Need to integrate with api TODO"

	return health, nil
}

func (s *grpcConsciousServer) MakeLLMChatRequest(prompt string) (string, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(prompt))
	completions, err := MakeLLMChatRequest(prompt, s.memory.Messages, s.memory.Model, s.memory.Temperature)
	if err != nil {
		return "", err
	}

	// add response to chat
	aiMsg := completions.Choices[0].Message.Content
	fmt.Println("AI:", aiMsg)
	s.memory.Messages = append(s.memory.Messages, openai.AssistantMessage(aiMsg))

	return aiMsg, nil

}
