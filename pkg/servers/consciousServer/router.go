package consciousServer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	protoV1 "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func (s *grpcConsciousServer) RouteRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("Router %s, %s", in.Type, in.Prompt)
	switch in.Type {
	case protoV1.ContextMesssage_LLM_QUERY:
		return s.RouteLLMChatRequest(in)
	case protoV1.ContextMesssage_LLM_IMAGE_DESCRIBE:
		return s.RouteLLMImageDescribeRequest(in)
	case protoV1.ContextMesssage_IMAGE_GENERATION:
		// use open AI API supporting generation
		return s.RouteOpenAiImageGenerationRequest(in)
		// return s.RouteImageCreateRequest(in)
	default:
		return nil, nil
	}

}

func (s *grpcConsciousServer) RouteLLMChatRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	s.memory.AppendMemory(openai.UserMessage(in.Prompt))
	completions, err := llm.MakeLLMChatRequest(s.memory.Messages, s.memory.Tools, s.memory.Model, s.memory.Temperature)
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
		s.memory.AppendMemory(openai.AssistantMessage(assistantMsgContent))
		s.memory.Save()

		// send to tts, for it to generate voice over of the text
		// mqttServer.SendOnTTS(assistantMsgContent)

	}

	return &protoV1.ContextMesssage{
		Uuid:        in.Uuid,
		Model:       in.Model,
		Prompt:      in.Prompt,
		CreatedAt:   in.CreatedAt,
		CompletedAt: utils.GetTimeStamp(),
		Data:        completions.Choices[0].Message.Content,
	}, nil

}

func (s *grpcConsciousServer) RouteLLMImageDescribeRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(in.Prompt))

	// load the image if BodyText is emtpy and data is present
	// then upload file is passed in data field
	if in.BodyText == "" && in.Data != "" {
		imageFileName := utils.GetUploadsFolder() + in.Data
		image, err := os.ReadFile(imageFileName)
		if err != nil {
			log.Fatalf("Failed to read image file: %v", err)
		}
		in.BodyText = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(image)
	}

	result, err := llm.MakeImageDescribeRequest(in.Prompt, in.SystemPrompt, in.Model, in.BodyText)
	if err != nil {
		return nil, err
	}

	assistantMsgContent := result.OutputText()
	s.memory.AppendMemory(openai.AssistantMessage(assistantMsgContent))

	return &protoV1.ContextMesssage{
		Uuid:        in.Uuid,
		Model:       in.Model,
		Prompt:      in.Prompt,
		CreatedAt:   in.CreatedAt,
		CompletedAt: utils.GetTimeStamp(),
		Response:    assistantMsgContent,
	}, nil

}

func (s *grpcConsciousServer) RouteOpenAiImageGenerationRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("RouteOpenAiImageGenerationRequest -> %s, %s", in.Model, in.Prompt)
	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(in.Prompt))

	result, err := llm.MakeOpenAiImageGenerationRequest(in.Prompt, in.Model)
	if err != nil {
		logrus.Errorf("RouteOpenAiImageGenerationRequest -> MakeOpenAiImageGenerationRequest result(%v), %v", result, err)
		return nil, err
	}

	logrus.Infof("RouteOpenAiImageGenerationRequest -> decoding image")
	image, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}

	FILE_STORE_PATH := utils.GetEnv("FILE_STORE_PATH", "./dataVolume")
	// servable file path from static server
	servableFile := in.Uuid + ".png"
	imageFileName := FILE_STORE_PATH + "/gen/" + servableFile

	logrus.Infof("RouteImageCreateRequest -> %s, %s, %s", in.Model, in.Prompt, imageFileName)
	if err := os.WriteFile(imageFileName, image, 0o600); err != nil {
		panic(err)
	}

	assistantMsgContent := "Image generated and saved to disk."
	s.memory.AppendMemory(openai.AssistantMessage(assistantMsgContent))

	return &protoV1.ContextMesssage{
		Uuid:          in.Uuid,
		Model:         in.Model,
		Prompt:        in.Prompt,
		CreatedAt:     in.CreatedAt,
		CompletedAt:   utils.GetTimeStamp(),
		ResponseImage: servableFile,
	}, nil

}

// func (s *grpcConsciousServer) RouteImageCreateRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

// 	s.memory.Messages = append(s.memory.Messages, openai.UserMessage(in.Prompt))

// 	FILE_STORE_PATH := utils.GetEnv("FILE_STORE_PATH", "./dataVolume")

// 	// servable file path from static server
// 	servableFile := in.Uuid + ".png"
// 	filename := FILE_STORE_PATH + "/gen/" + servableFile
// 	logrus.Infof("RouteImageCreateRequest -> %s, %s, %s", in.Model, in.Prompt, filename)
// 	// servableFile = "fa23d376-bc55-4ad6-94a0-dfc1ec04c6bf.png"
// 	err := llm.MakeImageGenerationRequest(in.Prompt, in.Model, filename)
// 	if err != nil {
// 		logrus.Errorf("RouteImageCreateRequest -> MakeImageGenerationRequest %v", err)
// 		return nil, err
// 	}

// 	return &protoV1.ContextMesssage{
// 		Uuid:          in.Uuid,
// 		Model:         in.Model,
// 		Prompt:        in.Prompt,
// 		CreatedAt:     in.CreatedAt,
// 		CompletedAt:   utils.GetTimeStamp(),
// 		ResponseImage: servableFile,
// 	}, nil

// }
