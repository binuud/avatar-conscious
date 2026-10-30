package memory

import (
	"fmt"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/settings"

	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	MEMORY_DIRECTORY = "memories"
)

type Memory struct {
	// extend the proto buffer
	pb.Memory
	Messages       []openai.ChatCompletionMessageParamUnion
	Tools          []openai.ChatCompletionToolUnionParam
	CurrentContext *pb.SingleContext
}

// NewMemory
// if name is provided, it will be persisted
// emtpy name will be used as scratch memory
// which will be removed after the session
func NewMemory(name ...string) *Memory {

	systemPrompt := "You are a helpful assistant. You must use the tools provided to answer questions. Very polite, friendly with children. Do not use curse words. Do not use bad words. Do not use any smilies, emojis, or symbols in your responses from now on."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	model := settings.Global.GetValue("LLM_MODEL")

	memoryName := ""
	if len(name) > 0 {
		memoryName = name[0]
	}

	memory := &Memory{
		Name:         memoryName,
		Model:        model,
		SystemPrompt: systemPrompt,
		Temperature:  0.8,
		Messages:     messages,
		Tools:        nil,
		Contexts:     make([]*pb.SingleContext, 0),
	}

	if memoryName == "" {
		// unnamed memory is use and throw, they are not saved
	} else {
		// if memory is named, try to load it from cache
		// memory with names _, is use and throw memory, ignore same
		memory.Load(memoryName)
	}

	return memory
}

func (m *Memory) DumpMemory(w io.Writer) {
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
		if w != nil {
			// write to an io writer, so UI client or non standard client
			// can view the output
			str := fmt.Sprintf("Name: %s, content: %s", role, content)
			w.Write([]byte(str))
		}
		logrus.Infof("Name: %s, content: %s", role, content)
	}
	logrus.Info("-----------------------------------------------------------------")
}

// AppendInMessage()
// Append IN Message, take the prompt and add to User Type message
// append to memory, only append user and assistant message to memory
func (m *Memory) AppendInMessage(in *pb.ContextMesssage) {
	// reset current context
	m.CurrentContext = &pb.SingleContext{
		In: in,
	}
	m.AppendMemory(openai.UserMessage(in.Prompt))
	m.Save()
}

// AppendOutMessage
// Append OUT Message, take the prompt and add to Assistant Type message
// append to memory, only append user and assistant message to memory
func (m *Memory) AppendOutMessage(out *pb.ContextMesssage) {
	if m.CurrentContext == nil {
		m.CurrentContext = &pb.SingleContext{}
		logrus.Errorf("AppendOutMessage: received an out message, without an in message")
	}
	m.CurrentContext.Out = out
	m.Contexts = append(m.Contexts, m.CurrentContext)
	m.AppendMemory(openai.AssistantMessage(out.Prompt))
	m.Save()
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
func (m *Memory) Load(name string) error {

	// Read profile back (Go automatically infers type T from your variable assignment)
	localMemory, err := datalayer.ReadBinary[Memory](datalayer.GetFileAbsolutePath(MEMORY_DIRECTORY, name))
	if err != nil {
		return err
	}
	*m = localMemory
	return nil
}

func (m *Memory) Save() error {

	if m.Name == "" {
		return nil
	}
	fullFileName := datalayer.GetFileAbsolutePath(MEMORY_DIRECTORY, m.Name)
	// Save list of configs
	if err := datalayer.SaveBinary(fullFileName, m); err != nil {
		logrus.Errorf("Cannot save file %s %v", fullFileName, err)
		return err
	}
	logrus.Infof("Saved memory to file %s", fullFileName)
	return nil
}
