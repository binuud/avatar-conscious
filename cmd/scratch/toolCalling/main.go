package main

import (
	"github.com/binuud/conscious/pkg/llm"
	"github.com/binuud/conscious/pkg/memory"
	apitools "github.com/binuud/conscious/pkg/tools/apiTools"
	"github.com/binuud/conscious/pkg/tools/commonTools"
	"github.com/binuud/conscious/pkg/tools/toolRegistry"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func CheckWebFetch() {
	memory := memory.NewMemory()

	tools := toolRegistry.NewToolRegistry()
	tools.AddTool(commonTools.NewToolFileContents())
	tools.AddTool(apitools.NewApiGetTool())

	memory.Tools = tools.GetCommonTools()

	prompt := "Get details about a person named john, using http get api request to url http://localhost:8081/user/[USER NAME]"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages, memory.Model, memory.Tools, memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.HandleToolCalls(msg.ToolCalls)
		if err == nil {
			logrus.Infof("Output received %s", out)
		}
	}
}

func CheckFileContents() {

	memory := memory.NewMemory()

	tools := toolRegistry.NewToolRegistry()
	tools.AddTool(commonTools.NewToolFileContents())
	tools.AddTool(apitools.NewApiGetTool())

	memory.Tools = tools.GetCommonTools()

	prompt := "What is the contents of AGENTS.md file"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages, memory.Model, memory.Tools, memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.HandleToolCalls(msg.ToolCalls)
		if err != nil {
			logrus.Errorf("Error received %v", err)
		} else {
			logrus.Infof("Output received %s", out)
		}

	}
}

func main() {
	// CheckFileContents()
	CheckWebFetch()
}
