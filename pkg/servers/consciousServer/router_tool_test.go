package consciousServer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/openai/openai-go/v3"
)

func Test_LLM_TOOL(t *testing.T) {

	memory := memory.NewMemory()

	tests := []struct {
		name       string
		prompt     string
		want_tool  string
		want_param string
		want_error bool
	}{
		{
			name:       "Check tools with weather api, city chennai",
			prompt:     "What is the current weather in Chennai",
			want_tool:  "get_weather",
			want_param: "Chennai",
			want_error: false,
		},
		{
			name:       "Check tools with weather api, city delhi",
			prompt:     "What is the current weather in Delhi",
			want_tool:  "get_weather",
			want_param: "Delhi",
			want_error: false,
		},
		{
			name:       "Check price of stock of Apple",
			prompt:     "What is the latest stock price of Apple",
			want_tool:  "get_price_of_stock",
			want_param: "Apple",
			want_error: false,
		},
		{
			name:       "Check price of stock of Google",
			prompt:     "Get me the stock price for Google",
			want_tool:  "get_price_of_stock",
			want_param: "Google",
			want_error: false,
		},
		{
			name:       "Image generation prompt",
			prompt:     "Create a realistic photo of a coffee cup, with steaming hot, and on a table",
			want_tool:  "get_image",
			want_param: "",
			want_error: false,
		},
		{
			name:       "Image generation prompt",
			prompt:     "Draw a image of a cattle grazing",
			want_tool:  "get_image",
			want_param: "",
			want_error: false,
		},
		{
			name:       "Image describe",
			prompt:     "Describe the different cars in the image",
			want_tool:  "get_image_description",
			want_param: "",
			want_error: false,
		},
		{
			name:       "Negative test case",
			prompt:     "What is your name",
			want_tool:  "",
			want_param: "",
			want_error: true,
		},
	}

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
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_price_of_stock",
			Description: openai.String("Get the price of a stock of a company or organization"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"companyName": map[string]string{
						"type":        "string",
						"description": "The name of the stock or company, eg: Apple, Google",
					},
				},
				"required": []string{"companyName"},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_image",
			Description: openai.String("Retreive or create an image"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"text": map[string]string{
						"type":        "string",
						"description": "A description of the image",
					},
				},
				"required": []string{"text"},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_image_description",
			Description: openai.String("Describe a given image or photo."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"text": map[string]string{
						"type":        "string",
						"description": "A question about the image",
					},
				},
				"required": []string{"text"},
			},
		}),
	}
	memory.Tools = tools

	// Loop through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			toolCall, err := routerViaTool(tt.prompt, memory, t)

			if err != nil {
				if tt.want_error {
					t.Logf("Test case failed as required %v", err)
				} else {
					t.Errorf("LLM Tool request errored: %v", err)
				}

			} else {

				if tt.want_error {
					t.Errorf("Expected error, but got none %v", err)
					return
				}

				t.Logf("Tool called : %s with arguments: %s\n", toolCall.Function.Name, toolCall.Function.Arguments)

				// DeepEqual safely compares arrays/slices element by element
				if !reflect.DeepEqual(toolCall.Function.Name, tt.want_tool) {
					t.Errorf("LLM Tooling failed, tool name not returned = %v, want %v", toolCall.Function.Name, tt.want_tool)
				}

				if toolCall.Function.Name == "get_weather" {
					var args struct {
						Location string `json:"location"`
					}
					json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

					// Mock result execution
					if !reflect.DeepEqual(args.Location, tt.want_param) {
						t.Errorf("LLM Tooling failed, tool argument not returned = %v, want %v", args.Location, tt.want_param)
					}

				}

				if toolCall.Function.Name == "get_price_of_stock" {
					var args struct {
						CompanyName string `json:"companyName"`
					}
					json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

					// Mock result execution
					if !reflect.DeepEqual(args.CompanyName, tt.want_param) {
						t.Errorf("LLM Tooling failed, tool argument not returned = %v, want %v", args.CompanyName, tt.want_param)
					}

				}

			}

		})
	}
}

func routerViaTool(prompt string, memory *memory.Memory, t *testing.T) (*openai.ChatCompletionMessageToolCallUnion, error) {

	t.Logf("LLM Prompt: %s", prompt)
	memory.AppendMemory(openai.UserMessage(prompt))
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages, memory.Model, memory.Tools, memory.Temperature)
	if err != nil {
		// t.Errorf("client: error llm chat request: %v \n", err)
		return nil, err
	}

	// logrus.Infof("Response for tool call %v", completions.Choices)

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {

		// The model wants to call a tool!
		toolCall := msg.ToolCalls[0]
		return &toolCall, nil

	}

	return nil, fmt.Errorf("no tools returned")

}
