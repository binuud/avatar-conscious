package main

import (
	"github.com/binuud/conscious/pkg/llm"
	"github.com/binuud/conscious/pkg/memory"
	"github.com/binuud/conscious/pkg/tools/commonTools"
	"github.com/binuud/conscious/pkg/tools/toolRegistry"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func main() {

	memory := memory.NewMemory()

	tools := toolRegistry.NewToolRegistry()
	tools.AddTool("get_file_contents", commonTools.NewToolFileContents())
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
		if err == nil {
			logrus.Infof("Output received %s", out)
		}
	}
}
