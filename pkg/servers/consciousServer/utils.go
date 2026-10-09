package consciousServer

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/sirupsen/logrus"
)

func GetLLMEndpoint() string {
	return os.Getenv("LLM_HOST") + "/v1/"
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
