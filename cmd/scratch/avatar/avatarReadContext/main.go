package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	apitools "github.com/binuud/avatar-conscious/pkg/tools/apiTools"
	"github.com/binuud/avatar-conscious/pkg/tools/commonTools"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

var tools = toolRegistry.NewToolRegistry()

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

func registerTools() {

	tools.AddTool(commonTools.NewToolFileContents())
	tools.AddTool(apitools.NewApiGetTool())
	tools.AddTool(aiGenTools.NewAiGenerateImageTool())
	tools.AddTool(aiGenTools.NewAiGenerateImageDescriptionTool())
}

func parseAvatar() {

	logrus.Infof("Read avatar file")
	memory := memory.NewMemory("")
	// Read the avatar.md file
	filePath := "dataVolume/avatarStore/election-officer.md"
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

	// step := map[string]any{
	// 	"type": "object",
	// 	"properties": map[string]any{
	// 		"name":        map[string]any{"type": "string"},
	// 		"description": map[string]any{"type": "string"},
	// 		"prompt":      map[string]any{"type": "string"},
	// 	},
	// 	"required":             []string{"explanation", "output"},
	// 	"additionalProperties": false,
	// }
	// schema := map[string]any{
	// 	"type": "object",
	// 	"properties": map[string]any{
	// 		"steps":       map[string]any{"type": "array", "items": step},
	// 		"name":        map[string]any{"type": "string"},
	// 		"description": map[string]any{"type": "string"},
	// 	},
	// 	"required":             []string{"steps", "final_answer"},
	// 	"additionalProperties": false,
	// }
	schema, err := getJsonSchemaMap()
	if err != nil {
		logrus.Fatalf("parseAvatar: %v", err)
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

func getJsonSchemaMap() (map[string]any, error) {

	var schemaMap map[string]any

	schemaJSON, err := jsonschema.For[Avatar](nil)
	if err != nil {
		log.Fatalf("Failed to run  jsonschema: %v", err)
	}

	output, err := json.MarshalIndent(schemaJSON, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal schema: %v", err)
	}

	// 3. Unmarshal the JSON bytes into your map
	err = json.Unmarshal(output, &schemaMap)
	if err != nil {
		fmt.Printf("Error unmarshalling schema: %v\n", err)
		return nil, err
	}
	return schemaMap, nil
}

func runAvatar() {
	output := `{
  "steps": [
    {
      "name": "Extract License Information",
      "description": "Analyze the uploaded license image to extract the license number and driver name.",
      "prompt": "Describe the image from file sampleLicense.jpg, and get the license number, remove all spaces from the license number, and driver name, in format {'license': [LICENCE NUMBER],'name':'[NAME OF PERSON]'}"
    },
    {
      "name": "Verify License Validity",
      "description": "Verify the extracted license number against the database API to confirm the user is a valid driver.",
      "prompt": "Call the API to verify the license number and compare the returned name with the extracted name. API Endpoint: http://localhost:8081/checklicense/[LICENSE NUMBER]"
    }
  ],
  "description": "Verifies if the user is a valid voter by checking the voter id with a database",
  "name": "election-officer"
	}`

	var in Avatar
	if err := json.Unmarshal([]byte(output), &in); err != nil {
		logrus.Fatalf("Cannot unmarshall %v", err)
	}

	memory := memory.NewMemory("")
	memory.Tools = tools.GetCommonTools()

	for i, step := range in.Steps {
		logrus.Infof("Step %d : %s", i, step.Name)
		logrus.Infof("Model : %s", memory.Model)
		logrus.Infof("Prompt %s", step.Prompt)
		memory.Messages = append(memory.Messages, openai.UserMessage(step.Prompt))

		completions, err := llm.MakeLLMSingleToolRequest(memory.Messages, memory.Model, memory.Tools, memory.Temperature)
		if err != nil {
			logrus.Errorf("client: error llm chat request: %v \n", err)
			return
		}

		msg := completions.Choices[0].Message
		if len(msg.ToolCalls) > 0 {
			out, err := tools.HandleToolCalls(msg.ToolCalls)
			if err != nil {
				logrus.Errorf("HandleToolCalls: %v \n", err)
				return
			}
			memory.Messages = append(memory.Messages, openai.AssistantMessage(out))
		}
	}

}

func main() {
	registerTools()

	parseAvatar()
	// runAvatar()
}
