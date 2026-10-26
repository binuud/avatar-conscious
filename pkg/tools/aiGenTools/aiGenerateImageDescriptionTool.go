package aiGenTools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// LLM returns a json string
// which has to be parsed into a equivalent structure
// this provides the input structure
type AiGenerateImageDescriptionTool_Input struct {
	// the uploaded or reference image file
	ImageFile string `json:"imagefile"`
	// the original prompt given by the user
	Prompt string `json:"prompt"`
}

// implements the CommonToolInterface
type AiGenerateImageDescriptionTool struct {
	toolMeta.CommonToolInterface
}

func NewAiGenerateImageDescriptionTool() *AiGenerateImageDescriptionTool {
	return &AiGenerateImageDescriptionTool{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t AiGenerateImageDescriptionTool) Name() string {
	return "describe_image"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t AiGenerateImageDescriptionTool) Description() string {
	return "Describe the given image or photograph, and collect some data from it."
}

func (t AiGenerateImageDescriptionTool) InputSchema() string {
	return "Describe the given image or photograph, and collect some data from it."
}

func (t AiGenerateImageDescriptionTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"imagefile": map[string]string{
				"type":        "string",
				"description": "The file which has to be described using ai",
			},
			"prompt": map[string]string{
				"type":        "string",
				"description": "The original prompt given by the user, do not change the user prompt, with instructions on return object",
			},
		},
		"required": []string{"imagefile", "prompt"},
	}
}

func (t AiGenerateImageDescriptionTool) Handler(in *pb.ContextMesssage) (*pb.ContextMesssage, error) {
	logrus.Infof("AiGenerateImageTool Handler Arguments received %v", in.Prompt)
	var inpObj AiGenerateImageDescriptionTool_Input
	if err := json.Unmarshal([]byte(in.Prompt), &inpObj); err != nil {
		return nil, err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)
	visionModel := utils.GetEnv("VISION_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit")
	outMsg, err := DescribeImage(inpObj.Prompt, visionModel, inpObj.ImageFile)
	if err != nil {
		return nil, err
	}
	return &pb.ContextMesssage{
		Uuid:     in.Uuid,
		Prompt:   in.Prompt,
		Model:    visionModel,
		Response: outMsg,
	}, nil

}

func DescribeImage(prompt string, model string, fileName string) (string, error) {

	memory := memory.NewMemory("")
	memory.Model = model
	memory.AppendMemory(openai.UserMessage(prompt))

	logrus.Infof("DescribeImage -> %s, %s", memory.Model, prompt)

	imageFileName := utils.GetUploadsFolder() + fileName
	image, err := os.ReadFile(imageFileName)
	if err != nil {
		return "", fmt.Errorf("Failed to read image file: %v", err)
	}
	base64Data := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(image)

	result, err := llm.MakeImageDescribeRequest(prompt, "", memory.Model, base64Data)
	if err != nil {
		return "", err
	}
	logrus.Infof("DescribeImage Output %s", result.OutputText())
	return result.OutputText(), nil

}
