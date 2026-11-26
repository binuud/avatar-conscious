package testTools

import (
	"testing"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/sirupsen/logrus"
)

func Test_UiLaunchTool(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bootstrap.Bootstrap()
	memoryName := ""
	memory := memory.NewMemory(memoryName)
	memory.InitTools()
	// memory.Tools = tools.GlobalCommonTools
	memory.AppendUserMessage("what is 10 + 10")
	memory.AppendAssistantMessage("20")
	memory.AppendUserMessage("Launch an Avatar UI interface, the ui interface id is : web-camera")
	_test_tool_calling(t, memory)

	memory.AppendUserMessage("what is 10 + 30")
	memory.AppendAssistantMessage("40")
	memory.AppendUserMessage("Step No: 1, image is stored in filename ```6bf47703-a591-46ec-9a18-01ea7dec713a.png```")
	memory.AppendAssistantMessage("ok")
	memory.AppendUserMessage("Launch Avatar UI interface with id : get-prompt")
	//memory.AppendMessage(openai.UserMessage("Step No: 2, Launch an Avatar UI interface, the ui interface id is : get-prompt"))
	memory.DumpMemory(nil)
	_test_tool_calling(t, memory)
}

func _test_tool_calling(t *testing.T, memory *memory.Memory) {
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Meta.Model,
		memory.Tools,
		memory.Meta.Temperature)
	if err != nil {
		t.Fatalf("client: error llm chat request: %v \n", err)
	}

	// check if any tool has to be called
	// if so call hanldle the tool and return the output response
	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		logrus.Infof("Toolcall %v", msg.ToolCalls)
		// memory.AppendMessage(openai.ToolMessage(msg.ToolCalls[0].RawJSON(), msg.ToolCalls[0].ID))
		out, err := memory.ToolRegistry.HandleToolCalls(msg.ToolCalls[0], "Step No: 2, Launch an Avatar UI interface, the ui interface id is : get-prompt")
		if err != nil {
			t.Errorf("HandleToolCalls: %v \n", err)
		} else {
			logrus.Infof("Tool %s called, Out Msg : %s", msg.ToolCalls[0].Function.Name, out)
			memory.AppendToolMessage("Avatar UI Interface with id web-camera shown to the user by the tool function", msg.ToolCalls[0].ID)
			// handle special case ui launch tool
			if msg.ToolCalls[0].Function.Name == "get_ui_launch_tool" {
				t.Log("Test tool call was made to get_ui_launch_tool")
			}
		}
	} else {
		t.Errorf("uiLaunchTool not called: %v \n", msg)
	}
}
