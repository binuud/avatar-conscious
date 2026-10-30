package aiGenTools

import (
	"encoding/base64"
	"encoding/json"
	"os"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// LLM returns a json string
// which has to be parsed into a equivalent structure
// this provides the input structure
type AiGenerateImageTool_Input struct {
	// refined text prompt
	Text string `json:"text"`
}

// implements the CommonToolInterface
type AiGenerateImageTool struct {
	toolMeta.CommonToolInterface
}

func NewAiGenerateImageTool() *AiGenerateImageTool {
	return &AiGenerateImageTool{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t AiGenerateImageTool) Name() string {
	return "get_ai_generated_image"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t AiGenerateImageTool) Description() string {
	return "Create or Modify an image using ai."
}

func (t AiGenerateImageTool) InputSchema() string {
	return "Create or Modify an image or photograph using ai."
}

func (t AiGenerateImageTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"text": map[string]string{
				"type":        "string",
				"description": "A refined and amplified description of the image or photograph",
			},
		},
		"required": []string{"text"},
	}
}

func (t AiGenerateImageTool) Handler(in *pb.ContextMesssage) (*pb.ContextMesssage, error) {
	logrus.Infof("AiGenerateImageTool Handler Arguments received %v", in.Prompt)
	var inpObj AiGenerateImageTool_Input
	if err := json.Unmarshal([]byte(in.Prompt), &inpObj); err != nil {
		return nil, err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)

	imageModel := settings.Global.GetValue("IMAGE_MODEL")

	outMsg, err := GetAiImage(inpObj.Text, imageModel)
	if err != nil {
		return nil, err
	}
	return &pb.ContextMesssage{
		Uuid:          in.Uuid,
		Prompt:        in.Prompt,
		Model:         imageModel,
		ResponseImage: outMsg,
	}, nil
}

func GetAiImage(prompt string, model string) (string, error) {

	memory := memory.NewMemory("")
	memory.Model = model
	memory.AppendMemory(openai.UserMessage(prompt))

	logrus.Infof("GetAiImage -> %s, %s", memory.Model, prompt)

	result, err := llm.MakeOpenAiImageGenerationRequest(prompt, memory.Model)
	if err != nil {
		return "", err
	}

	logrus.Infof("GetAiImage -> decoding image")
	image, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}

	SERVABLE_DIRECTORY := settings.Global.GetValue("SERVABLE_DIRECTORY")
	// servable file path from static server
	servableFile := "/gen/images/" + utils.GenerateUUID() + ".png"
	imageFileName := SERVABLE_DIRECTORY + servableFile

	logrus.Infof("GetAiImage -> %s, %s, %s", memory.Model, prompt, imageFileName)
	if err := os.WriteFile(imageFileName, image, 0o600); err != nil {
		panic(err)
	}
	return servableFile, nil
}
