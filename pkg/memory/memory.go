package memory

import (
	datalayer "github.com/binuud/conscious/pkg/dataLayer"
	"github.com/binuud/conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	MEMORY_DIRECTORY = "memories"
)

type Memory struct {
	Name         string
	Model        string
	SystemPrompt string
	Temperature  float64
	Messages     []openai.ChatCompletionMessageParamUnion
	Tools        []openai.ChatCompletionToolUnionParam
}

func NewMemory() *Memory {

	systemPrompt := "You are a helpful assistant. You must use the tools provided to answer questions. Very polite, friendly with children. Do not use curse words. Do not use bad words. Do not use any smilies, emojis, or symbols in your responses from now on."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	model := utils.GetEnv("LLM_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit")

	memory := &Memory{
		Name:         "test-memory",
		Model:        model,
		SystemPrompt: systemPrompt,
		Temperature:  0.8,
		Messages:     messages,
		Tools:        nil,
	}
	return memory
}

func (m *Memory) DumpMemory() {
	logrus.Info("-----------------------------------------------------------------")
	logrus.Infof("Name: %s, Model: %s, Temp: %f", m.Name, m.Model, m.Temperature)
	logrus.Infof("System Prompt: %s", m.SystemPrompt)
	for _, msg := range m.Messages {
		var role string
		var content string
		switch {
		case msg.OfUser != nil:
			role = "User"
			content = msg.OfUser.Content.OfString.Value
		case msg.OfAssistant != nil:
			role = "Assistant"
			content = msg.OfAssistant.Content.OfString.Value
		case msg.OfSystem != nil:
			role = "System"
			content = msg.OfSystem.Content.OfString.Value
		case msg.OfDeveloper != nil:
			role = "Developer"
			content = msg.OfDeveloper.Content.OfString.Value
		case msg.OfFunction != nil:
			role = "Function"
			content = msg.OfFunction.Content.Value
		case msg.OfTool != nil:
			role = "Tool"
			content = msg.OfTool.Content.OfString.Value
		}

		logrus.Infof("Name: %s, content: %s", role, content)
	}
	logrus.Info("-----------------------------------------------------------------")
}

// AppendMemory
// append to memory, only append user and assistant message to memory
func (m *Memory) AppendMemory(message openai.ChatCompletionMessageParamUnion) {
	if message.OfUser != nil || message.OfAssistant != nil {
		m.Messages = append(m.Messages, message)
	}
	m.Save()
}

// Load function to load from model file
func (m *Memory) Load() error {

	// Read profile back (Go automatically infers type T from your variable assignment)
	localMemory, err := datalayer.ReadBinary[Memory](datalayer.GetFileAbsolutePath(MEMORY_DIRECTORY, m.Name))
	if err != nil {
		return err
	}
	*m = localMemory
	return nil
}

func (m *Memory) Save() error {

	fullFileName := datalayer.GetFileAbsolutePath(MEMORY_DIRECTORY, m.Name)
	// Save list of configs
	if err := datalayer.SaveBinary(fullFileName, m); err != nil {
		logrus.Errorf("Cannot save file %s %v", fullFileName, err)
		return err
	}
	logrus.Infof("Saved memory to file %s", fullFileName)
	return nil
}
