package consciousServer

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

func MakeImageDescribeRequest(prompt string, model string) (*responses.Response, error) {

	llmURL := os.Getenv("LLM_ENDPOINT") + "/engines/llama.cpp/v1/"
	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentParamOfInputText(prompt),
						{OfInputImage: &responses.ResponseInputImageParam{
							Detail:   responses.ResponseInputImageDetailAuto,
							ImageURL: openai.String("https://api.nga.gov/iiif/a2e6da57-3cd1-4235-b20e-95dcaefed6c8/full/!800,800/0/default.jpg"),
						}},
					},
					responses.EasyInputMessageRoleUser,
				),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(response.OutputText())
	return response, nil

}

func MakeLLMChatRequest(prompt string, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolParam, model string, temperature float64) (*openai.ChatCompletion, error) {

	llmURL := os.Getenv("LLM_ENDPOINT") + "/engines/llama.cpp/v1/"

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
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

	llmURL := os.Getenv("LLM_ENDPOINT") + "/engines/llama.cpp/v1/"

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
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
