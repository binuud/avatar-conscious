package llm

import (
	"context"
	"fmt"
	"io"
	"os"

	protoV1 "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/settings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
	"github.com/openai/openai-go/v3/responses"

	"github.com/sirupsen/logrus"
)

func GetLLMEndpoint() string {
	llmHost := settings.Global.GetValue("LLM_ENDPOINT")
	if llmHost == "" {
		logrus.Fatalf("llm endpoint is not set, please check env LLM_HOST variable")
	}
	return llmHost + "/v1/"
}

func GetLLMClient() openai.Client {
	return openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(settings.Global.GetValue("LLM_ENDPOINT")), // No API key needed for Model Runner
	)
}

func MakeLLMGetModels() ([]*protoV1.Model, error) {

	client := GetLLMClient()
	// var models []protoV1.Model
	models := make([]*protoV1.Model, 0)

	page, err := client.Models.List(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to llm server %v", err)
	}

	logrus.Println("Available Models:")
	// Loop through the results automatically handled by the SDK iterator
	for _, item := range page.Data {
		model := &protoV1.Model{
			Id:      item.ID,
			Uuid:    "",
			Created: item.Created,
			OwnedBy: item.OwnedBy,
		}
		logrus.Printf("- %s (Owned by: %s)\n", model.Id, model.OwnedBy)
		models = append(models, model)
	}

	return models, nil

}

func MakeImageDescribeRequest(prompt string, systemPrompt string, model string, imageDataUrl string) (*responses.Response, error) {

	client := GetLLMClient()
	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentParamOfInputText(prompt),
						{OfInputImage: &responses.ResponseInputImageParam{
							Detail:   responses.ResponseInputImageDetailAuto,
							ImageURL: openai.String(imageDataUrl),
						}},
					},
					responses.EasyInputMessageRoleUser,
				),
			},
		},
	})

	if err != nil {
		logrus.Errorf("panic MakeImageDescribeRequest %s", err)
		panic(err)
	}

	fmt.Println(response.OutputText())
	return response, nil
}

// MakeLLMChatStreamingChatRequest
// Create a streaming request, this returns a stream
// this requests also maintians history
func MakeLLMChatStreamingChatRequest(messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam, model string, temperature float64) *ssestream.Stream[openai.ChatCompletionChunk] {

	client := GetLLMClient()

	ctx := context.Background()

	// logrus.Infof("MakeLLMChatRequest :%s:%s", model, prompt)
	param := openai.ChatCompletionNewParams{
		Messages:        messages,
		Model:           model,
		Temperature:     openai.Opt(temperature),
		Tools:           tools,
		ReasoningEffort: "none",
	}

	// Create the streaming request
	stream := client.Chat.Completions.NewStreaming(ctx, param)

	return stream

}

// MakeLLMChatRequest
// Create a NON streaming request with history, this returns a completed response
// this requests also maintians history
func MakeLLMChatRequest(messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam, model string, temperature float64) (*openai.ChatCompletion, error) {

	client := GetLLMClient()

	ctx := context.Background()

	// logrus.Infof("MakeLLMChatRequest :%s:%s", model, prompt)
	param := openai.ChatCompletionNewParams{
		Messages:        messages,
		Model:           model,
		Temperature:     openai.Opt(temperature),
		Tools:           tools,
		ReasoningEffort: "none",
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		return nil, err
	}

	return completion, nil

}

// MakeLLMSingleRequest
// Create a NON streaming request with out history, this returns a completed response
// this requests DOES NOT maintian history
func MakeLLMSingleRequest(prompt string, systemPrompt string, model string, temperature float64) (string, error) {

	client := GetLLMClient()

	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(prompt),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(temperature),
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		return "", err
	}

	var response = completion.Choices[0].Message.Content

	fmt.Println(response)
	return response, nil

}

func MakeLLMSchemaResponseRequest(prompt string, systemPrompt string, model string, schema map[string]any) (string, error) {

	client := GetLLMClient()

	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{OfInputItemList: responses.ResponseInputParam{
			responses.ResponseInputItemParamOfMessage(
				responses.ResponseInputMessageContentListParam{responses.ResponseInputContentParamOfInputText(systemPrompt)},
				responses.EasyInputMessageRoleSystem,
			),
			responses.ResponseInputItemParamOfMessage(
				responses.ResponseInputMessageContentListParam{responses.ResponseInputContentParamOfInputText(prompt)},
				responses.EasyInputMessageRoleUser,
			),
		}},
		Text: responses.ResponseTextConfigParam{Format: responses.ResponseFormatTextConfigUnionParam{
			OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{Name: "event", Schema: schema, Strict: openai.Bool(true)},
		}},
	})
	if err != nil {
		panic(err)
	}

	return response.OutputText(), nil

}

func MakeLLMSingleToolRequest(messages []openai.ChatCompletionMessageParamUnion, model string, tools []openai.ChatCompletionToolUnionParam, temperature float64) (*openai.ChatCompletion, error) {

	client := GetLLMClient()

	ctx := context.Background()

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Tools:       tools,
		Temperature: openai.Opt(temperature),
	}

	return client.Chat.Completions.New(ctx, param)

}

// GetLLMToolCall
// returns the tool name, that the model wants to call
func GetLLMToolName(messages []openai.ChatCompletionMessageParamUnion, model string, tools []openai.ChatCompletionToolUnionParam, temperature float64) (*openai.ChatCompletionMessageToolCallUnion, error) {

	completions, err := MakeLLMSingleToolRequest(messages, model, tools, temperature)
	if err != nil {
		// t.Errorf("client: error llm chat request: %v \n", err)
		return nil, err
	}
	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {

		// The model wants to call a tool!
		toolCall := msg.ToolCalls[0]
		return &toolCall, nil

	}
	return nil, nil
}

func MakeOpenAiImageGenerationRequest(prompt string, model string) (*openai.ImagesResponse, error) {

	client := GetLLMClient()

	return client.Images.Generate(context.Background(), openai.ImageGenerateParams{
		Model:          openai.ImageModel(model),
		Prompt:         prompt,
		N:              openai.Int(1),
		ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON,
		Size:           openai.ImageGenerateParamsSize1024x1024,
	})

}

func MakeOpenAiImageEditRequest(prompt string, model string, imageFileName string) (*openai.ImagesResponse, error) {

	client := GetLLMClient()

	sampleFile, err := os.Open(imageFileName)
	if err != nil {
		return nil, fmt.Errorf("Error reading sample file: %v\n", err)
	}
	defer sampleFile.Close()

	return client.Images.Edit(context.Background(), openai.ImageEditParams{
		Image:          openai.ImageEditParamsImageUnion{OfFile: io.Reader(sampleFile)},
		Model:          openai.ImageModel(model),
		Prompt:         prompt,
		N:              openai.Int(1),
		ResponseFormat: openai.ImageEditParamsResponseFormatB64JSON,
		Size:           openai.ImageEditParamsSize1024x1024,
	})

}
