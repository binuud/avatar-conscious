package commonTools

// Binu Udayakumar
// binu@dronasys.com
// Do not change this file without my approval

// Tool - fileContents
// Get the contents of a file, whose path is given by LLM
// This can be used as reference to generate other others tools

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
type Tool_FileContents_Input struct {
	// path is the absolute path to read (required).
	Path string `json:"path"`
}

// implements the CommonToolInterface
type ToolFileContents struct {
	toolMeta.CommonToolInterface
}

func NewToolFileContents() *ToolFileContents {
	return &ToolFileContents{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t ToolFileContents) Name() string {
	return "get_file_contents"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
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
