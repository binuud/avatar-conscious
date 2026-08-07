package consciousServer

import (
	"context"
	"fmt"
	"os"

	protoV1 "github.com/binuud/conscious/gen/go/v1/brain"
	aiModuleProto "github.com/dronasys-com/brain-nodes/gen/go/v1/aiModule"
	mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/sirupsen/logrus"
)

var mqttClient mqtt.Client

func InitMQTT() {

	// MQTT_HOST Mqtt Server with port
	MQTT_HOST := os.Getenv("MQTT_HOST")
	if MQTT_HOST == "" {
		logrus.Panic("env MQTT_HOST not set")
	}

	opts := mqtt.NewClientOptions()

	opts.AddBroker(MQTT_HOST) // Set your MQTT broker
	opts.SetClientID("go-websocket-mqtt-client")

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		logrus.Fatalf("MQTT connect error: %v", token.Error())
	}
}

type grpcConsciousServer struct {
	protoV1.UnimplementedConsciousServer
	model        string
	systemPrompt string
	temperature  float64
	messages     []openai.ChatCompletionMessageParamUnion
}

func NewGrpcConsciousServer() *grpcConsciousServer {

	systemPrompt := "You are a very helpful assistant. Very polite, friendly with children. Do not use curse words. Do not use bad words. Reply briefly, in less than 20 words as far as possible. Do not use any smilies, emojis, or symbols in your responses from now on."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	// InitMQTT()

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
	mqttServer.SendOnTTS(response)

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
