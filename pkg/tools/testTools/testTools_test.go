package testTools

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/openai/openai-go/v3"
)

func Test_Tool_Registry(t *testing.T) {
	bootstrap.Bootstrap()
	memoryName := ""
	memory := memory.NewMemory(memoryName)
	memory.Tools = tools.GlobalCommonTools

	tests := []struct {
		name       string
		prompt     string
		want_tool  string
		want_param string
		want_error bool
	}{
		// {
		// 	name:       "Check tools with weather api, city chennai",
		// 	prompt:     "What is the current weather in Chennai",
		// 	want_tool:  "get_weather",
		// 	want_param: "Chennai",
		// 	want_error: false,
		// },
		// {
		// 	name:       "Check tools with weather api, city delhi",
		// 	prompt:     "What is the current weather in Delhi",
		// 	want_tool:  "get_weather",
		// 	want_param: "Delhi",
		// 	want_error: false,
		// },
		// {
		// 	name:       "Check price of stock of Apple",
		// 	prompt:     "What is the latest stock price of Apple",
		// 	want_tool:  "get_price_of_stock",
		// 	want_param: "Apple",
		// 	want_error: false,
		// },
		// {
		// 	name:       "Check price of stock of Google",
		// 	prompt:     "Get me the stock price for Google",
		// 	want_tool:  "get_price_of_stock",
		// 	want_param: "Google",
		// 	want_error: false,
		// },
		{
			name:       "Image generation prompt",
			prompt:     "Create a realistic photo of a coffee cup, with steaming hot, and on a table",
			want_tool:  "get_ai_generated_image",
			want_param: "",
			want_error: false,
		},
		{
			name:       "Image generation prompt",
			prompt:     "Draw a image of a cattle grazing",
			want_tool:  "get_ai_generated_image",
			want_param: "",
			want_error: false,
		},
		{
			name:       "Image describe",
			prompt:     "Describe the different cars in the image",
			want_tool:  "describe_image",
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
		{
			name:       "UILaunchTool",
			prompt:     "Launch an UI on Avatar to get some details",
			want_tool:  "get_ui_launch_tool",
			want_param: "",
			want_error: false,
		},
	}

	// Loop through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			toolCall, err := _getLLMToolName(tt.prompt, memory, t)

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

			}

		})
	}
}

func _getLLMToolName(prompt string, memory *memory.Memory, t *testing.T) (*openai.ChatCompletionMessageToolCallUnion, error) {

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}
	completions, err := llm.MakeLLMSingleToolRequest(messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		t.Fatalf("client: error llm chat request: %v \n", err)
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		return &msg.ToolCalls[0], nil
	}

	return nil, fmt.Errorf("no tools returned")

}
