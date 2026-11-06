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
type AiGenerateImageEditTool_Input struct {
	// the uploaded or reference image file
	ImageFiles []string `json:"imagefiles"`
	// the original prompt given by the user
	Prompt string `json:"prompt"`
}

// implements the CommonToolInterface
type AiGenerateImageEditTool struct {
	toolMeta.CommonToolInterface
}

func NewAiGenerateImageEditTool() *AiGenerateImageEditTool {
	return &AiGenerateImageEditTool{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t AiGenerateImageEditTool) Name() string {
	return "get_ai_edit_image"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t AiGenerateImageEditTool) Description() string {
	return "Edit the given image, or use the given image as reference to create a new image"
}

func (t AiGenerateImageEditTool) InputSchema() string {
	return "Edit the given image, or use the given image as reference to create a new image"
}

func (t AiGenerateImageEditTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"imagefiles": map[string]any{
				"type":        "array",
				"description": "The image file names with extension png, jpg for get_ai_edit_image",
				"items": map[string]any{
					"type": "string",
				},
			},
			"prompt": map[string]any{
				"type":        "string",
				"description": "A refined and amplified instructions for editing the image or photograph",
			},
		},
		"required": []string{"imagefiles", "prompt"}, // This must be at the same level as "properties"
	}
}

func (t AiGenerateImageEditTool) Handler(in *pb.ContextMesssage) (*pb.ContextMesssage, error) {
	logrus.Infof("AiGenerateImageEditTool Handler Arguments received %v", in.Prompt)
	var inpObj AiGenerateImageEditTool_Input
	if err := json.Unmarshal([]byte(in.Prompt), &inpObj); err != nil {
		return nil, err
	}
	logrus.Infof("AiGenerateImageEditTool: Arguments parsed into input structure %v", inpObj)

	imageModel := settings.Global.GetValue("IMAGE_MODEL")

	// hotswap original prompt
	outMsg, err := GetAiImageEdit(&inpObj, imageModel)
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

func GetAiImageEdit(inpObj *AiGenerateImageEditTool_Input, model string) (string, error) {

	memory := memory.NewMemory("")
	memory.Model = model
	memory.AppendMemory(openai.UserMessage(inpObj.Prompt))
	logrus.Infof("GetAiImageEdit -> %s, %s, %s", memory.Model, inpObj.ImageFiles, inpObj.Prompt)

	imageFileNames := make([]string, 0)
	for _, fileName := range inpObj.ImageFiles {
		imageFileNames = append(imageFileNames, utils.GetUploadsFolder()+fileName)
	}

	result, err := llm.MakeOpenAiImageEditRequest(inpObj.Prompt, memory.Model, imageFileNames)
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
	newImageFileName := SERVABLE_DIRECTORY + servableFile

	logrus.Infof("GetAiImageEdit saving final image-> %s, %s, %s", memory.Model, inpObj.Prompt, newImageFileName)
	if err := os.WriteFile(newImageFileName, image, 0o600); err != nil {
		panic(err)
	}
	return servableFile, nil
}
