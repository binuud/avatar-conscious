package consciousServer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/binuud/conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
)

func createMemoryObject() *memory {

	systemPrompt := "You are a helpful weather assistant. You must use the tools provided to answer questions. Very polite, friendly with children. Do not use curse words. Do not use bad words. Reply briefly, in less than 20 words as far as possible. Do not use any smilies, emojis, or symbols in your responses from now on."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	model := utils.GetEnv("LLM_MODEL", "gemma4:latest")

	memory := &memory{
		Name:         "Default",
		Model:        model,
		SystemPrompt: systemPrompt,
		Temperature:  0.8,
		Messages:     messages,
		Tools:        nil,
	}
	return memory
}

func TestConsciousServer_LLMChat(t *testing.T) {

	memory := createMemoryObject()

	prompt := "What is 8 + 8"
	t.Logf("LLM Prompt: %s", prompt)
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	completions, err := MakeLLMChatRequest(prompt, memory.Messages, memory.Tools, memory.Model, memory.Temperature)
	if err != nil {
		t.Errorf("client: error llm chat request: %v \n", err)
	}

	assistantMsgContent := completions.Choices[0].Message.Content
	t.Logf("Received Response: %s", assistantMsgContent)

}

func TestConsciousServer_LLMImageProcess(t *testing.T) {
	// does not work
	// image describe request does not work
	memory := createMemoryObject()

	// // memory.Model = "ai/gemma3:4B-Q4_K_M"

	fileStore := os.Getenv("FILE_STORE_PATH")
	if fileStore == "" {
		t.Fatal("fatal error: FILE_STORE_PATH env not set")
	}

	imagePath := fileStore + "sampleImages/elephantWithBall.jpg" // Replace with your image file path
	t.Logf("Image file for testing %s", imagePath)
	image, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}
	imageDataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)

	response, err := MakeImageDescribeRequest("What's in this image? ", memory.Model, imageDataURL)
	if err != nil {
		t.Errorf("client: error llm image request: %v \n", err)
	} else {
		expectedSubstring := "elephant"
		t.Logf("Response image request %s", response.OutputText())
		if !strings.Contains(response.OutputText(), expectedSubstring) {
			t.Errorf("expected output to contain %q", expectedSubstring)
		}

	}

	if err != nil {
		panic(err)
	}

	fmt.Println(response.OutputText())

}

// func TestConsciousServer_LLMImageProcess(t *testing.T) {

// 	fileStore := os.Getenv("FILE_STORE_PATH")
// 	if fileStore == "" {
// 		t.Fatal("fatal error: FILE_STORE_PATH env not set")
// 	}

// 	memory := createMemoryObject()
// 	// 2. Read the local image file into memory
// 	imagePath := fileStore + "/sample.jpg" // Replace with your image file path
// 	imageBytes, err := os.ReadFile(imagePath)
// 	if err != nil {
// 		log.Fatalf("Failed to read image file: %v", err)
// 	}
// 	// 4. Encode image bytes to a Base64 string
// 	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
// 	mimeType := http.DetectContentType(imageBytes)
// 	// 5. Construct the data URL format required by OpenAI
// 	imageDataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)

// 	prompt := "Analyze this image and describe what you see in detail."
// 	t.Logf("LLM Prompt: %s", prompt)

// 	newMessage := []openai.ChatCompletionMessageParamUnion{
// 		openai.UserMessage(&openai.ChatCompletionUserMessageParam{
// 			Content: &openai.ChatCompletionUserMessageParamContentUnion{},
// 			Role:    "User",
// 		}),
// 		// openai.ImageContentPart(imageDataURL),
// 	}

// 	memory.Messages = newMessage
// 	completions, err := MakeLLMChatRequest(prompt, memory.Messages, memory.Tools, memory.Model, memory.Temperature)
// 	if err != nil {
// 		t.Errorf("client: error llm chat request: %v \n", err)
// 	}

// 	assistantMsgContent := completions.Choices[0].Message.Content
// 	t.Logf("Received Response: %s", assistantMsgContent)

// }

func TestConsciousServer_LLMToolFunction(t *testing.T) {

	memory := createMemoryObject()

	// test case for checking tools,
	// test with a sample weather function
	tools := []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_weather",
			Description: openai.String("Get the current weather for a given city or place"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]string{
						"type":        "string",
						"description": "The city or place, e.g. Tokyo",
					},
				},
				"required": []string{"location"},
			},
		}),
	}
	memory.Tools = tools

	prompt := "What is the current weather in Chennai City"
	t.Logf("LLM Prompt: %s", prompt)
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	completions, err := MakeLLMChatRequest(prompt, memory.Messages, memory.Tools, memory.Model, memory.Temperature)
	if err != nil {
		t.Errorf("client: error llm chat request: %v \n", err)
	}

	assistantMsgContent := ""
	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {

		// The model wants to call a tool!
		toolCall := msg.ToolCalls[0]
		t.Logf("Model invoked tool: %s with arguments: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)

		// Parse arguments and run your local go function logic here
		if toolCall.Function.Name == "get_weather" {
			var args struct {
				Location string `json:"location"`
			}
			json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			// Mock result execution
			if args.Location == "Chennai" {
				t.Error("client: error llm tool works, but argument is not correctly sent my llm\n")
			} else {
				t.Logf(" %s City, information correctly sent as argument\n", args.Location)
			}
			assistantMsgContent = fmt.Sprintf("Weather in %s, is Sunny, and the temperature is 35C", args.Location)

			t.Logf("Tools args %v", args)
			t.Logf("Received Tool Response: %s", assistantMsgContent)
		}

	} else {
		t.Error("client: error Tool Chain not called\n")
	}

}
