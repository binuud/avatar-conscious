package main

import (
	"encoding/base64"
	"os"

	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func CheckImageDescriber() string {

	memory := memory.NewMemory("")
	memory.Tools = tools.GlobalCommonTools

	prompt := "Describe the image from file sampleLicense.jpg, and get the license number and driver name, in format {'license': [LICENSE NUMBER],'name':'[NAME OF PERSON']}"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))

	logrus.Infof("LLM Prompt: %s", prompt)

	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return ""
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return ""
		}
		logrus.Infof("Output received %s", out)
		return out.Response
	}
	return msg.Content
}

func CheckImageGeneration() string {

	memory := memory.NewMemory("")
	memory.Tools = tools.GlobalCommonTools

	prompt := "Create an image of a cat, crawling on a desk, trying to catch a mouse"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))

	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return ""
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls)
		if err == nil {
			logrus.Infof("Output received %s", out)
			return out.Response
		}
	}
	return msg.Content
}

func CheckImageUpgradeGeneration() {
	memory := memory.NewMemory("")
	memory.Tools = tools.GlobalCommonTools

	prompt := "Create an image of a tourist card, use the photo of the person from image samplePerson.jpg as reference, show the person in front of eiffel tower"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))

	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		logrus.Fatalf("client: error llm chat request: %v \n", err)

	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls)
		if err == nil {
			logrus.Infof("Output received %s", out)

		} else {
			logrus.Fatalf("client: tool call: %v \n", err)
		}
	}

}

func CheckDirectImageUpgradeGeneration() {

	fileName := "samplePerson.jpg"
	imageFileName := utils.GetUploadsFolder() + fileName

	imageModel := utils.GetEnv("IMAGE_MODEL", "Runpod/FLUX.2-klein-4B-mflux-4bit")
	prompt := "Create an image of a tourist card, use the photo of the person from image samplePerson.png as reference, show the person in front of eiffel tower"
	result, err := llm.MakeOpenAiImageEditRequest(prompt, imageModel, imageFileName)
	if err != nil {
		logrus.Fatalf("CheckDirectImageUpgradeGeneration err, %v", err)
	}

	logrus.Infof("GetAiImage -> decoding image")
	image, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}

	SERVABLE_DIRECTORY := utils.GetEnv("SERVABLE_DIRECTORY", "./dataVolume/static")
	// servable file path from static server
	servableFile := "/gen/images/" + utils.GenerateUUID() + ".png"
	newImageFileName := SERVABLE_DIRECTORY + servableFile

	logrus.Infof("GetAiImage -> %s, %s, %s", imageModel, prompt, newImageFileName)
	if err := os.WriteFile(newImageFileName, image, 0o600); err != nil {
		panic(err)
	}

}

func CheckWebFetch(systemPrompt string) string {
	memory := memory.NewMemory("")
	memory.Tools = tools.GlobalCommonTools

	prompt := "Get license details of the person, strip spaces in the license number , using http get api request to url http://localhost:8081/checklicense/[license Number]"
	memory.Messages = append(memory.Messages, openai.AssistantMessage("License details are as follows "+systemPrompt))
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return ""
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls)
		if err == nil {
			logrus.Infof("Output received %s", out)
			return out.Response
		}
	}
	return msg.Content
}

func CheckFileContents() string {

	memory := memory.NewMemory("")
	memory.Tools = tools.GlobalCommonTools

	prompt := "What is the contents of AGENTS.md file"
	memory.Messages = append(memory.Messages, openai.UserMessage(prompt))
	logrus.Infof("LLM Prompt: %s", prompt)
	completions, err := llm.MakeLLMSingleToolRequest(memory.Messages,
		memory.Model,
		memory.Tools,
		memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return ""
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls)
		if err == nil {
			logrus.Infof("Output received %s", out)
			return out.Response
		}
	}
	return msg.Content
}

func main() {
	bootstrap.Bootstrap()
	// out := CheckImageDescriber()
	// if out == "" {
	// 	logrus.Fatalf("No output received from CheckImageDescriber")
	// }
	// we got output from previous step
	// CheckWebFetch(out)
	// CheckImageGeneration()
	CheckImageUpgradeGeneration()
	// CheckDirectImageUpgradeGeneration()
	// CheckFileContents()

}
