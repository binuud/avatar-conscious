package testTools

import (
	"testing"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func Test_ImageEdit(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bootstrap.Bootstrap()
	memoryName := ""
	memory := memory.NewMemory(memoryName)
	memory.Tools = tools.GlobalCommonTools
	memory.AppendMemory(openai.UserMessage("what is 10 + 10"))
	memory.AppendMemory(openai.AssistantMessage("20"))
	memory.AppendMemory(openai.UserMessage(`Edit the image "e9420877-7407-4c79-95ee-932065a99ac0.png", 
		use the dress provided in image "shirt.jpg"
		using the image of the person from webcamera as reference, 
		Do not modify the face of the person.
		Keep the facial features, skin tone of the person exactly as they are in the original image.
		create a new image of the person wearing the dress, do not change face or skin tone`))

	var expectedTool = "get_ai_edit_image"
	var expectedParams = []string{"e9420877-7407-4c79-95ee-932065a99ac0.png", "shirt.jpg"}
	_test_tool_request_params(t, memory, expectedTool, expectedParams)

	memory.AppendMemory(openai.AssistantMessage("Image was generated"))
	memory.AppendMemory(openai.UserMessage(`Edit the image "e9420877-7407-4c79-95ee-932065a99ac0.png", 
		using the image of the person from webcamera as reference, 
		Do not modify the face of the person.
		Keep the facial features, skin tone of the person exactly as they are in the original image.
		create a new image of the person wearing the dress, do not change face or skin tone`))

	_test_tool_request_params(t, memory, expectedTool, expectedParams)
}

func _test_tool_request_params(t *testing.T, memory *memory.Memory, expectedTool string, expectedParams []string) {
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		t.Fatalf("client: error llm chat request: %v \n", err)
	}

	// check if any tool has to be called
	// if so call hanldle the tool and return the output response
	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		toolCall := msg.ToolCalls[0]
		logrus.Infof("Function %+v", toolCall.Function.Name)
		logrus.Infof("Args %+v", toolCall.Function.Arguments)
		if expectedTool == toolCall.Function.Name {
			t.Logf("tool name matches expected %s", expectedTool)
			// now check params
		} else {
			t.Errorf("tool name does not match expected %s", expectedTool)
		}

	} else {
		t.Errorf("uiLaunchTool not called: %v \n", msg)
	}
}
