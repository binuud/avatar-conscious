package consciousServer

// Model Router
// Handles memory on per request basis, allowing users to
// use different tabs with chats connecting to different models

import (
	"encoding/json"
	"fmt"
	"strings"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	protoV1 "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

type router struct {
	memory *memory.Memory
	in     *pb.ContextMesssage
}

func NewRouter(in *pb.ContextMesssage) *router {
	memory := memory.NewMemory(in.MemoryName)
	in.Model = memory.Meta.Model
	in.Seed = memory.Meta.Seed
	in.Temperature = float32(memory.Meta.Temperature)
	logrus.Infof("Initialized new Router with memory %s", in.MemoryName)
	return &router{
		memory: memory,
		in:     in,
	}
}

func (r *router) RouteStreamingPromptIntent(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	out := &pb.ContextMesssage{}
	var err error

	switch utils.DetectImageIntent(in.Prompt) {
	// case utils.IntentDescribe:
	// 	in.Model = settings.Global.GetValue("VISION_MODEL")
	// 	out, err = r.RouteLLMImageDescribeRequest(in)
	case utils.IntentCreate:
		in.Model = settings.Global.GetValue("IMAGE_MODEL")
		out = r.CreateOutStreamMessage(in, fmt.Sprintf("Creating Image with model: %s\n", in.Model))
		if err := stream.Send(out); err != nil {
			return err
		}
		out, err = r.RouteOpenAiImageGenerationRequest(in)
	// case utils.IntentEdit:
	// 	in.Model = settings.Global.GetValue("IMAGE_MODEL")
	// 	out, err = r.RouteOpenAiImageEditRequest(in)
	default:
		return r.RouteStreamWithToolsRequest(in, stream)
	}
	if err = stream.Send(out); err != nil {
		return err
	}
	return nil
}

func (r *router) RouteRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("Router %s, %s", in.Type, in.Prompt)
	switch in.Type {
	case protoV1.ContextMesssage_LLM_QUERY:
		return r.RouteLLMChatRequest(in)
	case protoV1.ContextMesssage_LLM_IMAGE_DESCRIBE:
		return r.RouteLLMImageDescribeRequest(in)
	case protoV1.ContextMesssage_IMAGE_GENERATION:
		// use open AI API supporting generation
		return r.RouteOpenAiImageGenerationRequest(in)
	case protoV1.ContextMesssage_IMAGE_EDIT:
		// use open AI API supporting generation
		return r.RouteOpenAiImageEditRequest(in)
	default:
		return nil, nil
	}

}

func (r *router) RouteLLMChatRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	r.memory.AppendInMessage(in)
	completions, err := llm.MakeLLMChatRequest(r.memory.Messages, r.memory.Tools, r.memory.Meta.Model, r.memory.Meta.Temperature)
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
			out := r.CreateAndAppendOutMessage(in, assistantMsgContent)
			return out, nil
		}
	}

	logrus.Println("Model direct response:", assistantMsgContent)
	out := r.CreateAndAppendOutMessage(in, assistantMsgContent)
	return out, nil

}

func (r *router) RouteLLMImageDescribeRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	r.memory.AppendInMessage(in)

	assistantMsgContent, err := aiGenTools.DescribeImage(in.Prompt, in.Model, in.Data)
	if err != nil {
		logrus.Errorf("RouteLLMImageDescribeRequest -> DescribeImage result(%v), %v", assistantMsgContent, err)
		return nil, err
	}

	logrus.Infof("RouteLLMImageDescribeRequest -> %s, %s, %s", in.Model, in.Prompt, assistantMsgContent)
	out := r.CreateAndAppendOutMessage(in, assistantMsgContent)

	return out, nil

}

func (r *router) RouteOpenAiImageGenerationRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("RouteOpenAiImageGenerationRequest -> %s, %s", in.Model, in.Prompt)

	r.memory.AppendInMessage(in)

	imageFileName, err := aiGenTools.GetAiImage(in.Prompt, in.Model)
	if err != nil {
		logrus.Errorf("RouteOpenAiImageGenerationRequest -> GetAiImage result(%s), %v", imageFileName, err)
		return nil, err
	}

	logrus.Infof("RouteImageCreateRequest -> %s, %s, %s", in.Model, in.Prompt, imageFileName)
	r.memory.AppendToolMessage(imageFileName, "")
	assistantMsgContent := "Image generated and saved to file: " + imageFileName
	out := r.CreateOutMessage(in, assistantMsgContent)
	out.ResponseImage = imageFileName
	r.memory.AppendOutMessage(out)

	return out, nil

}

func (r *router) RouteOpenAiImageEditRequest(in *protoV1.ContextMesssage) (*protoV1.ContextMesssage, error) {

	logrus.Infof("RouteOpenAiImageEditRequest -> %s, %s", in.Model, in.Prompt)
	r.memory.AppendInMessage(in)

	inpObj := &aiGenTools.AiGenerateImageEditTool_Input{
		Prompt:     in.Prompt,
		ImageFiles: strings.Split(in.Data, ","),
	}

	imageFileName, err := aiGenTools.GetAiImageEdit(inpObj, in.Model)
	if err != nil {
		logrus.Errorf("RouteOpenAiImageEditRequest -> GetAiImageEdit result(%v), %v", imageFileName, err)
		return nil, err
	}

	logrus.Infof("RouteOpenAiImageEditRequest -> %s, %s, %s", in.Model, in.Prompt, imageFileName)

	assistantMsgContent := "Image edited and saved to file: " + imageFileName
	out := r.CreateOutMessage(in, assistantMsgContent)
	out.ResponseImage = imageFileName
	r.memory.AppendOutMessage(out)

	return out, nil

}

func (r *router) RouteStreamRequest(in *protoV1.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	// thinking := ""
	completeResponse := ""

	r.memory.AppendInMessage(in)

	llmStream := llm.MakeLLMChatStreamingChatRequest(r.memory.Messages, r.memory.Tools, r.memory.Meta.Model, r.memory.Meta.Temperature, r.memory.Meta.Seed)

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
	out := r.CreateAndAppendOutMessage(in, completeResponse)

	// send the out message
	if err := stream.Send(out); err != nil {
		return err
	}

	return nil
}

func (r *router) RouteStreamWithToolsRequest(in *protoV1.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) error {

	// thinking := ""
	completeResponse := ""

	r.memory.AppendInMessage(in)

	for {
		llmStream := llm.MakeLLMChatStreamingChatRequest(r.memory.Messages, r.memory.Tools, r.memory.Meta.Model, r.memory.Meta.Temperature, r.memory.Meta.Seed)

		defer func() {
			_ = llmStream.Close()
		}()

		logrus.Infof("InStream Received -> Request: %s, %s", in.Type, in.Prompt)
		acc := openai.ChatCompletionAccumulator{}
		// Iterate through the stream chunks as they arrive
		for llmStream.Next() {

			chunk := llmStream.Current()
			acc.AddChunk(chunk)
			if len(chunk.Choices) > 0 {
				// Print each chunk to stdout in real-time
				fmt.Print(chunk.Choices[0].Delta.Content)
				out := &pb.ContextMesssage{
					Uuid:     in.Uuid,
					Response: chunk.Choices[0].Delta.Content,
				}
				completeResponse = completeResponse + chunk.Choices[0].Delta.Content
				// 	// send the out message
				if err := stream.Send(out); err != nil {
					return err
				}
			}
			// Optional: trigger execution as soon as an individual tool finishes streaming
			if tool, ok := acc.JustFinishedToolCall(); ok {
				fmt.Printf("\n[Tool ready: %s(%s)]\n", tool.Name, tool.Arguments)
			}

		}
		if err := llmStream.Err(); err != nil {
			logrus.Errorf("stream error %v", err)
			return err
		}

		// Check accumulated message for tool calls
		msg := acc.Choices[0].Message
		if len(msg.ToolCalls) == 0 {
			logrus.Info("Stream closed")
			break
		}

		// Append the assistant's tool-call message to history
		r.memory.AppendDirectMessage(msg.ToParam())
		// params.Messages = append(params.Messages, msg.ToParam())

		// Execute each tool and append its result
		for _, tc := range msg.ToolCalls {
			out := r.CreateOutMessage(in, fmt.Sprintf("Calling tool %s (%s)\n", tc.Function.Name, tc.Function.Arguments))
			if err := stream.Send(out); err != nil {
				return err
			}
			out, err := r.memory.ToolRegistry.HandleToolCalls(tc, in.Prompt)
			if err != nil {
				logrus.Errorf("HandleToolCalls: %v \n", err)
				// return err
			}
			toolMessage := out.Response
			if out.ResponseImage != "" {
				toolMessage = "Image generated and saved to file: " + out.ResponseImage
				out.Response = toolMessage

				r.memory.AppendOutMessage(out)
				return nil
			} else {
				r.memory.AppendToolMessage(toolMessage, tc.ID)
			}
			if err := stream.Send(out); err != nil {
				return err
			}
		}

	}

	// if streaming is done
	// send last out message with completeAt timestamp
	out := r.CreateAndAppendOutMessage(in, completeResponse)

	// send the out message
	if err := stream.Send(out); err != nil {
		return err
	}

	return nil
}

func (r *router) MakeLLMChatRequest(prompt string) (*openai.ChatCompletion, error) {

	r.memory.AppendUserMessage(prompt)
	completions, err := llm.MakeLLMChatRequest(r.memory.Messages, r.memory.Tools, r.memory.Meta.Model, r.memory.Meta.Temperature)
	if err != nil {
		return nil, err
	}

	return completions, nil

}

// CreateAndAppendOutMessage
// Created and clones an out Context message from an in message
// sets the assistant content, and appends to memory
func (r *router) CreateAndAppendOutMessage(in *pb.ContextMesssage, assistantMessage string) *pb.ContextMesssage {
	out := r.CreateOutMessage(in, assistantMessage)
	r.memory.AppendOutMessage(out)
	return out
}

func (r *router) CreateOutMessage(in *pb.ContextMesssage, assistantMessage string) *pb.ContextMesssage {
	out := &protoV1.ContextMesssage{
		Uuid:        in.Uuid,
		Model:       in.Model,
		Prompt:      in.Prompt,
		CreatedAt:   in.CreatedAt,
		CompletedAt: utils.GetTimeStamp(),
		Response:    assistantMessage,
	}
	return out
}

func (r *router) CreateOutStreamMessage(in *pb.ContextMesssage, assistantMessage string) *pb.ContextMesssage {
	out := &protoV1.ContextMesssage{
		Uuid:      in.Uuid,
		Model:     in.Model,
		Prompt:    in.Prompt,
		CreatedAt: in.CreatedAt,
		Response:  assistantMessage,
	}
	return out
}
