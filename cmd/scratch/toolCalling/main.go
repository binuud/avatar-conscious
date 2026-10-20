package main

import (
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	apitools "github.com/binuud/avatar-conscious/pkg/tools/apiTools"
	"github.com/binuud/avatar-conscious/pkg/tools/commonTools"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

var tools = toolRegistry.NewToolRegistry()

func registerTools() {

	tools.AddTool(commonTools.NewToolFileContents())
	tools.AddTool(apitools.NewApiGetTool())
	tools.AddTool(aiGenTools.NewAiGenerateImageTool())
	tools.AddTool(aiGenTools.NewAiGenerateImageDescriptionTool())
}

func CheckImageDescriber() string {

	memory := memory.NewMemory()
	memory.Tools = tools.GetCommonTools()
	prompt := "Describe the image from file sampleLicense.jpg, and get the license number and driver name, in format {'license': [LICENSE NUMBER],'name':'[NAME OF PERSON']}"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))

	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages, memory.Model, memory.Tools, memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return ""
	}
	out := ""
	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err = tools.HandleToolCalls(msg.ToolCalls)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return ""
		}

	}
	return out
}

func CheckImageGeneration() {

	memory := memory.NewMemory()
	memory.Tools = tools.GetCommonTools()

	prompt := "Create an image of a cat, crawling on a desk, trying to catch a mouse"
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

func CheckWebFetch(systemPrompt string) {
	memory := memory.NewMemory()
	memory.Tools = tools.GetCommonTools()

	prompt := "Get license details of the person, strip spaces in the license number , using http get api request to url http://localhost:8081/checklicense/[license Number]"
	memory.Messages = append(memory.Messages, openai.AssistantMessage("License details are as follows "+systemPrompt))
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
			logrus.Errorf("CheckWebFetch: %v \n", err)
		}

		logrus.Infof("Output received %s", out)

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
	registerTools()
	out := CheckImageDescriber()
	if out == "" {
		logrus.Fatalf("No output received from CheckImageDescriber")
	}
	// we got output from previous step
	CheckWebFetch(out)
	// CheckImageGeneration()
	// CheckFileContents()

}
