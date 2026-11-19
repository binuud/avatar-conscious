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
		out, err = r.RouteOpenAiImageGenerationRequest(in)
	// case utils.IntentEdit:
	// 	in.Model = settings.Global.GetValue("IMAGE_MODEL")
	// 	out, err = r.RouteOpenAiImageEditRequest(in)
	default:
		// for tool calling, do not send the history
		_, out, err = r.handleTool(in, stream)
		if err != nil {
			return err
		}
		// special case, when there is an image generation,
		// we do not want llm to handle it
		if out != nil && out.ResponseImage != "" {
			r.memory.AppendOutMessage(out)
			return nil
		}
		// this is a normal chat streaming call
		// user has not asked for any tools
		return r.RouteStreamRequest(in, stream)
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
	// return &protoV1.ContextMesssage{
	// 	Uuid:        in.Uuid,
	// 	Model:       in.Model,
	// 	Prompt:      in.Prompt,
	// 	CreatedAt:   in.CreatedAt,
	// 	CompletedAt: utils.GetTimeStamp(),
	// 	Data:        completions.Choices[0].Message.Content,
	// }, nil

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

	llmStream := llm.MakeLLMChatStreamingChatRequest(r.memory.Messages, r.memory.Tools, r.memory.Meta.Model, r.memory.Meta.Temperature)

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

// handleTool
// returns a bool, and error. The boolean is to indicate whether the tool was handled
func (r *router) handleTool(in *pb.ContextMesssage, stream grpc.ServerStreamingServer[pb.ContextMesssage]) (bool, *pb.ContextMesssage, error) {

	// adding user prompt, so that the tooling call can decide what tool to use
	r.memory.AppendInMessage(in)
	completions, err := llm.MakeLLMSingleToolRequest(r.memory.Messages,
		r.memory.Meta.Model,
		r.memory.Tools,
		r.memory.Meta.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return false, nil, err
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {

		out, err := r.memory.ToolRegistry.HandleToolCalls(msg.ToolCalls, in.Prompt)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return false, nil, err
		}
		toolMessage := out.Response
		if out.ResponseImage != "" {
			toolMessage = out.ResponseImage
		}
		r.memory.AppendToolMessage(toolMessage, msg.ToolCalls[0].ID)
		// for metrics, populate the time stamps
		out.CreatedAt = in.CreatedAt
		out.CompletedAt = utils.GetTimeStamp()

		logrus.Debugf("Generate Out Msg : %s", out)
		// s.memory.AppendOutMessage(out)
		if err := stream.Send(out); err != nil {
			return true, out, err
		}
		return true, out, nil
	} else {

		// logrus.Debugf("Generate Tool function NOT selected\n")
		// logrus.Debugf("Generate Out Msg : %s", msg.Content)
		// if in.Data != "" {
		// 	// if an image is in the attachment, do a describe image
		// 	in.Type = pb.ContextMesssage_LLM_IMAGE_DESCRIBE
		// 	in.Model = settings.Global.GetValue("VISION_MODEL")
		// }
	}
	return false, nil, nil
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
