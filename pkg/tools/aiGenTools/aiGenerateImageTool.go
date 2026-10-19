package aiGenTools

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/binuud/conscious/pkg/llm"
	"github.com/binuud/conscious/pkg/memory"
	"github.com/binuud/conscious/pkg/tools/toolMeta"
	"github.com/binuud/conscious/pkg/utils"
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

func (t AiGenerateImageTool) Handler(args string) (string, error) {
	logrus.Infof("Arguments received %v", args)
	var in AiGenerateImageTool_Input
	if err := json.Unmarshal([]byte(args), &in); err != nil {
		return "", err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)
	return GetAiImage(in.Text)
}

func GetAiImage(prompt string) (string, error) {

	memory := memory.NewMemory()
	memory.Model = utils.GetEnv("IMAGE_MODEL", "Runpod/FLUX.2-klein-4B-mflux-4bit")
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

	FILE_STORE_PATH := utils.GetEnv("FILE_STORE_PATH", "./dataVolume")
	// servable file path from static server
	servableFile := utils.GenerateUUID() + ".png"
	imageFileName := FILE_STORE_PATH + "/gen/" + servableFile

	logrus.Infof("GetAiImage -> %s, %s, %s", memory.Model, prompt, imageFileName)
	if err := os.WriteFile(imageFileName, image, 0o600); err != nil {
		panic(err)
	}
	return servableFile, nil
}
