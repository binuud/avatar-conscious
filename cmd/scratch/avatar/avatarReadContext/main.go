package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/binuud/conscious/pkg/llm"
	"github.com/binuud/conscious/pkg/memory"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

type Step struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

type Avatar struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
}

func parseAvatar() {

	logrus.Infof("Read avatar file")
	memory := memory.NewMemory()
	// Read the avatar.md file
	filePath := "dataVolume/avatars/election-officer.md"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading %s: %v\n", filePath, err)
		return
	}
	avatarContent := string(fileBytes)

	// Formulate your prompt instructions
	systemPrompt := "Review the provided instruction structure and extract a JSON list of all steps with prompts if needed, api endpoints and commands to be run."

	completePrompt := fmt.Sprintf(
		"Here is the avatar context file (%s):\n\n```markdown\n%s\n```\n\n",
		filePath,
		avatarContent,
	)

	memory.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(memory.SystemPrompt),
		openai.UserMessage(completePrompt),
	}

	// logrus.Infof("skill prompt %s", completePrompt)

	// schema, err := getJSONSchema()
	// if err != nil {
	// 	logrus.Fatalf("client: error llm chat request: %v \n", err)
	// }

	step := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
			"prompt":      map[string]any{"type": "string"},
		},
		"required":             []string{"explanation", "output"},
		"additionalProperties": false,
	}
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"steps":       map[string]any{"type": "array", "items": step},
			"name":        map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
		},
		"required":             []string{"steps", "final_answer"},
		"additionalProperties": false,
	}

	output, err := llm.MakeLLMSchemaResponseRequest(completePrompt, systemPrompt, memory.Model, schema)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
	}

	var in Avatar
	if err := json.Unmarshal([]byte(output), &in); err != nil {
		logrus.Fatalf("Cannot unmarshall %v", err)
	}
	fmt.Println("\n✨ output:")
	fmt.Println(output)
	fmt.Printf("Struct %+v", in)
}

func getJSONSchema() (*jsonschema.Schema, error) {
	// Infer the schema from the Go struct type
	schema, err := jsonschema.For[Avatar](nil)
	if err != nil {
		log.Fatalf("Failed to run  jsonschema: %v", err)
	}
	return schema, nil

	// // Marshal to JSON for printing
	// output, err := json.MarshalIndent(schema, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Failed to marshal schema: %v", err)
	// }

	// fmt.Println(string(output))
	// return string(output), nil
}

func main() {
	parseAvatar()
}
