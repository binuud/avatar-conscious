package consciousServer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	protoV1 "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func (s *grpcConsciousServer) RouteRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("Router %s, %s", in.Type, in.Prompt)
	switch in.Type {
	case protoV1.ContextMesssage_LLM_QUERY:
		return s.RouteLLMChatRequest(in)
	case protoV1.ContextMesssage_IMAGE_GENERATION:
		return s.RouteLLMImageGenerationRequest(in)
	default:
		return nil, nil
	}

}

func (s *grpcConsciousServer) RouteLLMChatRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(in.Prompt))
	completions, err := MakeLLMChatRequest(in.Prompt, s.memory.Messages, s.memory.Tools, s.memory.Model, s.memory.Temperature)
	if err != nil {
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
		// mqttServer.SendOnTTS(assistantMsgContent)

	}

	return &protoV1.ContextMesssage{
		Data: completions.Choices[0].Message.Content,
	}, nil

}

func (s *grpcConsciousServer) RouteLLMImageGenerationRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("RouteLLMImageGenerationRequest -> %s, %s", in.Model, in.Prompt)
	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(in.Prompt))
	result, err := MakeLLMImageGenerationRequest(in.Prompt, in.Model)
	if err != nil {
		logrus.Errorf("RouteLLMImageGenerationRequest -> MakeLLMImageGenerationRequest result(%v), %v", result, err)
		return nil, err
	}
	logrus.Infof("RouteLLMImageGenerationRequest -> decoding image %v", result)
	image, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("/app/"+in.Uuid+".png", image, 0o600); err != nil {
		panic(err)
	}

	return &protoV1.ContextMesssage{
		Uuid:          in.Uuid,
		ResponseImage: "otter.png",
	}, nil

}
