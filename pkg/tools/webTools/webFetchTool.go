package webtools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// LLM returns a json string
// which has to be parsed into a equivalent structure
// this provides the input structure
type WebFetchTool_Input struct {
	// path is the absolute path to read (required).
	Url string `json:"url"`
}

// implements the CommonToolInterface
type WebFetchTool struct {
	toolMeta.CommonToolInterface
}

func NewWebFetchTool() *WebFetchTool {
	return &WebFetchTool{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t WebFetchTool) Name() string {
	return "get_file_contents"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t WebFetchTool) Description() string {
	return "Get the contents of the file, from the given path"
}

func (t WebFetchTool) InputSchema() string {
	return "Get the contents of the file, from the given path"
}

func (t WebFetchTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
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

func (t WebFetchTool) Handler(args string) (string, error) {
	logrus.Infof("Arguments received %v", args)
	var in WebFetchTool_Input
	if err := json.Unmarshal([]byte(args), &in); err != nil {
		return "", err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)
	return GetFileContents(in.Url)
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
