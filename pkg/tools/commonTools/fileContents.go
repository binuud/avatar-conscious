package commonTools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/binuud/conscious/pkg/tools/toolMeta"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

type Tool_FileContents_Input struct {
	// path is the absolute path to read (required).
	Path string `json:"path"`
}

type ToolFileContents struct {
	toolMeta.CommonToolInterface
}

func NewToolFileContents() *ToolFileContents {
	return &ToolFileContents{}
}

func (t ToolFileContents) Name() string {
	return "get_file_contents"
}

func (t ToolFileContents) Description() string {
	return "Get the contents of the file, from the given path"
}

func (t ToolFileContents) InputSchema() string {
	return "Get the contents of the file, from the given path"
}

func (t ToolFileContents) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]string{
				"type":        "string",
				"description": "The relative path of the file.",
			},
		},
		"required": []string{"path"},
	}
}

func (t ToolFileContents) Handler(args string) (string, error) {
	logrus.Infof("Arguments received %v", args)
	var in Tool_FileContents_Input
	if err := json.Unmarshal([]byte(args), &in); err != nil {
		return "", err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)
	return GetFileContents(in.Path)
}

func GetFileContents(fileName string) (string, error) {
	if strings.HasPrefix(fileName, "..") {
		return "", fmt.Errorf("security block: directory traversal detected")
	}
	// Read the entire file into a byte slice
	out, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
