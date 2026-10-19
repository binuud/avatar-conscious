package consciousServer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	pb "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/binuud/conscious/pkg/llm"
	"github.com/binuud/conscious/pkg/memory"
	"github.com/binuud/conscious/pkg/utils"
	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
)

func TestConsciousServer_LLMChat(t *testing.T) {

	memory := memory.NewMemory()

	prompt := "What is 8 + 8"
	t.Logf("LLM Prompt: %s", prompt)
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	completions, err := llm.MakeLLMChatRequest(memory.Messages, memory.Tools, memory.Model, memory.Temperature)
	if err != nil {
		t.Errorf("client: error llm chat request: %v \n", err)
	}

	assistantMsgContent := completions.Choices[0].Message.Content
	t.Logf("Received Response: %s", assistantMsgContent)

}

func TestConsciousServer_LLMImageProcess(t *testing.T) {
	// does not work
	// image describe request does not work
	memory := memory.NewMemory()

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
	imageDataURL := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(image)

	response, err := llm.MakeImageDescribeRequest("What's in this image? ", "answer in less than 10 words.", memory.Model, imageDataURL)
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

func TestConsciousServer_LLMToolFunction(t *testing.T) {

	memory := memory.NewMemory()

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
	completions, err := llm.MakeLLMChatRequest(memory.Messages, memory.Tools, memory.Model, memory.Temperature)
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

func TestConsciousServer_ImageGeneration(t *testing.T) {

	t.Logf("Testing Image creation, then validate same by using ai, to see if the image matches the prompt")
	s := NewGrpcConsciousServer()
	in := &pb.ContextMesssage{
		Uuid:   uuid.New().String(),
		Prompt: "Create an realistic phot of a cat throwing a ball to an elephant. The ball should have the words 'binuud'. The cat on the left side, and elephant on the right",
		Model:  "Runpod/FLUX.2-klein-4B-mflux-4bit",
		Type:   pb.ContextMesssage_IMAGE_GENERATION,
	}
	out, err := s.RouteRequest(in)
	if err != nil {
		t.Errorf("TestConsciousServer_ImageGeneration: cannot create image : %v \n", err)
	} else {
		t.Logf("Image generated and stored in (%v)\n", out.ResponseImage)

		fileStore := os.Getenv("FILE_STORE_PATH")
		if fileStore == "" {
			t.Fatal("fatal error: FILE_STORE_PATH env not set")
		}

		imagePath := fileStore + "/gen/" + out.ResponseImage
		image, err := os.ReadFile(imagePath)
		if err != nil {
			log.Fatalf("Failed to read image file: %v", err)
		}
		imageDataURL := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(image)

		in := &pb.ContextMesssage{
			Uuid:         uuid.New().String(),
			Prompt:       "What animals are present in this image, what is the text message on the ball ?, Reply in less than 10 words, do not output thinking, just mention the animals and text found in the image, do not use any comma, just space seperated list of words",
			Model:        "lmstudio-community/Qwen3.5-9B-MLX-4bit",
			Type:         pb.ContextMesssage_LLM_IMAGE_DESCRIBE,
			SystemPrompt: "Reply in less than 10 words, do not output thinking, just mention the animals and text found in the image",
			BodyText:     imageDataURL,
		}

		out, err := s.RouteRequest(in)
		if err != nil {
			t.Errorf("cannot describe the image : %v \n", err)
		} else {
			t.Logf("Image Describe works %s \n", out.Response)
			needleWords := [...]string{"elephant", "cat", "binuud"}
			if utils.ContainsAllWords(out.Response, needleWords[:]) {
				t.Logf("checked generated image with ai, to see if it contains prompted items, output contains all required words %v \n", needleWords)
			} else {
				t.Errorf("generated image might not contain the prompted items %v \n", needleWords)
			}
		}
	}
}
