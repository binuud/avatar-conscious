package main

import (
	"fmt"
	"os"

	"github.com/binuud/conscious/pkg/servers/consciousServer"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.Infof("Read skills file")
	memory := consciousServer.NewMemory()
	// Read the skills.md file
	filePath := "dataVolume/skills/skill.md"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading %s: %v\n", filePath, err)
		return
	}
	skillsContent := string(fileBytes)

	// Formulate your prompt instructions
	userPrompt := "Review the provided skills structure and extract a JSON list of all skill names."

	completePrompt := fmt.Sprintf(
		"Here is the context file (%s):\n\n```markdown\n%s\n```\n\nTask: %s",
		filePath,
		skillsContent,
		userPrompt,
	)

	memory.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(memory.SystemPrompt),
		openai.UserMessage(completePrompt),
	}
	// logrus.Infof("skill prompt %s", completePrompt)

	completions, err := consciousServer.MakeLLMChatRequest(memory.Messages, nil, memory.Model, memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
	}

	if err != nil {
		fmt.Printf("❌ API Error: %v\n", err)
		return
	}

	fmt.Println("\n✨ Response:")
	fmt.Println(completions.Choices[0].Message.Content)
}
