// consciousServer implements the ConsciousServer for the service from conscious proto file

// This server, listens on In endpoint, for user prompts, then makes a connection with backend AI,
// gets the response and returns same.
// When it gets a response from ai model, it publishes same to MQTT broker, so tts and other services, can listen to response events
// and respond for same
// States are saved, so each time the server is started, the old chat is maintained

package consciousServer

import (
	"context"
	"encoding/json"
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
	Tools        []openai.ChatCompletionToolParam
}

type grpcConsciousServer struct {
	protoV1.UnimplementedConsciousServer
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

	tools := []openai.ChatCompletionToolParam{
		{
			Type: "function",
			Function: openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("Get the current weather for a given city or place"),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "The city or place, e.g. Tokyo",
						},
					},
					"required":             []string{"location"},
					"additionalProperties": false,
				},
			},
		},
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

func (s *grpcConsciousServer) In(ctx context.Context, in *protoV1.InRequest) (*protoV1.InResponse, error) {

	logrus.Infof("Sending prompt: %s", in.Prompt)
	completions, err := s.MakeLLMChatRequest(in.Prompt)
	if err != nil {
		logrus.Errorf("LLM Error %v", err)
		return nil, err
	}

	msg := completions.Choices[0].Message
	assistantMsgContent := completions.Choices[0].Message.Content

	if len(msg.ToolCalls) > 0 {

		// The model wants to call a tool!
		toolCall := msg.ToolCalls[0]
		fmt.Printf("Model invoked tool: %s with arguments: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)

		// Parse arguments and run your local go function logic here
		if toolCall.Function.Name == "get_weather" {
			var args struct {
				Location string `json:"location"`
			}
			json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			// Mock result execution
			//assistantMsgContent := fmt.Sprintf(`{"location": "%s", "temp": "35C", "condition": "Sunny"}`, args.Location)
			assistantMsgContent = fmt.Sprintf("Weather in %s, is Sunny, and the temperature is 35C", args.Location)
			logrus.Infof("Tools args %v", args)
			logrus.Println("Tool execution result:", assistantMsgContent)
			s.memory.Messages = append(s.memory.Messages, openai.AssistantMessage(assistantMsgContent))
		}

	} else {

		// add prompt and response to memory
		// provided response is not a tool call and a direct assistant response
		logrus.Println("Model direct response:", assistantMsgContent)
		// add response to memory
		s.memory.Messages = append(s.memory.Messages, openai.AssistantMessage(assistantMsgContent))

		// send to tts, for it to generate voice over of the text
		mqttServer.SendOnTTS(assistantMsgContent)

	}

	err = s.memory.Save("default")
	if err != nil {
		logrus.Panicf("Cannot save conscious memory %v", err)
	}

	return &protoV1.InResponse{
		Uuid:     in.Uuid,
		Prompt:   in.Prompt,
		Response: assistantMsgContent,
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

func (s *grpcConsciousServer) MakeLLMChatRequest(prompt string) (*openai.ChatCompletion, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(prompt))
	completions, err := MakeLLMChatRequest(prompt, s.memory.Messages, s.memory.Tools, s.memory.Model, s.memory.Temperature)
	if err != nil {
		return nil, err
	}

	return completions, nil

}
