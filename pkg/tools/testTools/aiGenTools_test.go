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

func Test_UiLaunchTool(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bootstrap.Bootstrap()
	memoryName := ""
	memory := memory.NewMemory(memoryName)
	memory.Tools = tools.GlobalCommonTools
	memory.AppendMemory(openai.UserMessage("what is 10 + 10"))
	memory.AppendMemory(openai.AssistantMessage("20"))
	memory.AppendMemory(openai.UserMessage("Launch an Avatar UI interface, the ui interface id is : web-camera"))
	_test_tool_calling(t, memory)

	memory.AppendMemory(openai.UserMessage("what is 10 + 30"))
	memory.AppendMemory(openai.AssistantMessage("40"))
	memory.AppendMemory(openai.UserMessage("Step No: 1, image is stored in filename ```6bf47703-a591-46ec-9a18-01ea7dec713a.png```"))
	memory.AppendMemory(openai.AssistantMessage("ok"))
	memory.AppendMemory(openai.UserMessage("Launch Avatar UI interface with id : get-prompt"))
	//memory.AppendMemory(openai.UserMessage("Step No: 2, Launch an Avatar UI interface, the ui interface id is : get-prompt"))
	memory.DumpMemory(nil)
	_test_tool_calling(t, memory)
}

func _test_tool_calling(t *testing.T, memory *memory.Memory) {
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
		logrus.Infof("Toolcall %v", msg.ToolCalls)
		// memory.AppendMemory(openai.ToolMessage(msg.ToolCalls[0].RawJSON(), msg.ToolCalls[0].ID))
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls, "Step No: 2, Launch an Avatar UI interface, the ui interface id is : get-prompt")
		if err != nil {
			t.Errorf("HandleToolCalls: %v \n", err)
		} else {
			logrus.Infof("Tool %s called, Out Msg : %s", msg.ToolCalls[0].Function.Name, out)
			memory.AppendMemory(openai.ToolMessage("Avatar UI Interface with id web-camera shown to the user by the tool function", msg.ToolCalls[0].ID))
			// handle special case ui launch tool
			if msg.ToolCalls[0].Function.Name == "get_ui_launch_tool" {
				t.Log("Test tool call was made to get_ui_launch_tool")
			}
		}
	} else {
		t.Errorf("uiLaunchTool not called: %v \n", msg)
	}
}
