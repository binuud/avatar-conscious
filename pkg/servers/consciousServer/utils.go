package consciousServer

import (
	"context"
	"fmt"
	"io"
	"os"

	protoV1 "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/seasonjs/hf-hub/api"
	sd "github.com/seasonjs/stable-diffusion"
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

func MakeImageDescribeRequest(prompt string, systemPrompt string, model string, imageDataUrl string) (*responses.Response, error) {

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

func MakeLLMChatRequest(prompt string, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam, model string, temperature float64) (*openai.ChatCompletion, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

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

func MakeLLMSingleToolRequest(prompt string, systemPrompt string, model string, tools []openai.ChatCompletionToolUnionParam, temperature float64) (*openai.ChatCompletion, error) {

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
		Tools:       tools,
		Temperature: openai.Opt(temperature),
	}

	return client.Chat.Completions.New(ctx, param)

}

func MakeOpenAiImageGenerationRequest(prompt string, model string) (*openai.ImagesResponse, error) {

	client := openai.NewClient(
		option.WithBaseURL(GetLLMEndpoint()),
		option.WithAPIKey(""), // No API key needed for Model Runner
	)

	return client.Images.Generate(context.Background(), openai.ImageGenerateParams{
		Model:          openai.ImageModel(model),
		Prompt:         prompt,
		N:              openai.Int(1),
		ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON,
		Size:           openai.ImageGenerateParamsSize1024x1024,
	})

}

func MakeImageGenerationRequest(prompt string, model string, fileName string) error {

	options := sd.DefaultOptions

	imageModel, err := sd.NewAutoModel(options)
	if err != nil {
		print(err.Error())
		return err
	}
	defer imageModel.Close()

	hapi, err := api.NewApi()
	if err != nil {
		return err
	}

	modelPath, err := hapi.Model("justinpinkney/miniSD").Get("miniSD.ckpt")
	if err != nil {
		return err
	}

	logrus.Infof("MakeImageGenerationRequest -> loadfile %s", modelPath)
	err = imageModel.LoadFromFile(modelPath)
	if err != nil {
		return err
	}
	var writers []io.Writer
	filenames := []string{
		fileName,
	}
	for _, filename := range filenames {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		writers = append(writers, file)
	}
	logrus.Info("MakeImageGenerationRequest -> start inference")
	err = imageModel.Predict(prompt, sd.DefaultFullParams, writers)
	if err != nil {
		return err
	}
	return nil
}
