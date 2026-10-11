package consciousServer

import (
	"context"
	"fmt"
	"os"

	protoV1 "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/sirupsen/logrus"
)

func GetLLMEndpoint() string {
	return os.Getenv("LLM_HOST") + "/v1/"
}

func MakeLLMGetModels() ([]*protoV1.Model, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	// var models []protoV1.Model
	models := make([]*protoV1.Model, 0)

	page, err := client.Models.List(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Available Models:")
	// Loop through the results automatically handled by the SDK iterator
	for _, item := range page.Data {
		model := &protoV1.Model{
			Id:      item.ID,
			Uuid:    "",
			Created: item.Created,
			OwnedBy: item.OwnedBy,
		}
		fmt.Printf("- %s (Owned by: %s)\n", model.Id, model.OwnedBy)
		models = append(models, model)
	}

	return models, nil

}

func MakeImageDescribeRequest(prompt string, model string, imageDataUrl string) (*responses.Response, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentParamOfInputText("What's in this image?"),
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

func MakeLLMChatRequest(prompt string, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam, model string, temperature float64) (*openai.ChatCompletion, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	ctx := context.Background()

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Opt(temperature),
		Tools:       tools,
	}

	completion, err := client.Chat.Completions.New(ctx, param)

	if err != nil {
		return nil, err
	}

	return completion, nil

}

func MakeLLMSingleRequest(prompt string, systemPrompt string, model string, temperature float64) (string, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

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

func MakeLLMImageGenerationRequest(prompt string, model string) (*openai.ImagesResponse, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	return client.Images.Generate(context.Background(), openai.ImageGenerateParams{
		Model:          openai.ImageModel(model),
		Prompt:         prompt,
		N:              openai.Int(1),
		ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON,
		Size:           openai.ImageGenerateParamsSize256x256,
	})

}
