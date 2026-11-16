package memory

import (
	"fmt"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"

	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	MEMORY_DIRECTORY = "memories"
	SYSTEM_PROMPT    = `You are a helpful assistant. 
		You must use the tools provided to answer questions. 
		Very polite, friendly with children. Do not use curse words. 
		Do not use bad words. Do not use any smilies, emojis, or symbols in your responses from now on.
		If tool call response if a filename or a filename with path, give output as "\render [filename]"
		`
)

type Memory struct {
	// extend the proto buffer
	pb.Memory
	Messages       []openai.ChatCompletionMessageParamUnion
	ToolRegistry   *toolRegistry.ToolRegistry
	Tools          []openai.ChatCompletionToolUnionParam
	CurrentContext *pb.SingleContext
}

// NewScratchMemory
// scratch memory is not saved
func NewScratchMemory() *Memory {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(SYSTEM_PROMPT),
	}

	model := settings.Global.GetValue("LLM_MODEL")

	memory := &Memory{
		Model:        model,
		SystemPrompt: SYSTEM_PROMPT,
		Temperature:  0.8,
		Messages:     messages,
		Tools:        nil,
		Contexts:     make([]*pb.SingleContext, 0),
	}

	return memory
}

// NewMemory
// if name is provided, it will be persisted
// emtpy name will be used as scratch memory
// which will be removed after the session
// if existing memory is found, the same is loaded and returned
// else a new memory object is created
func NewMemory(name string) *Memory {

	memory, err := LoadMemory(name)
	if err != nil {
		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(SYSTEM_PROMPT),
		}

		model := settings.Global.GetValue("LLM_MODEL")

		memory = &Memory{
			Name:         name,
			Model:        model,
			SystemPrompt: SYSTEM_PROMPT,
			Temperature:  0.8,
			Messages:     messages,
			Tools:        nil,
			Contexts:     make([]*pb.SingleContext, 0),
		}
		memory.InitTools()
		memory.Save()
	}

	return memory
}

func (m *Memory) SetModel() {
	m.Model = settings.Global.GetValue("LLM_MODEL")
}

func (m *Memory) InitTools() {

	myTools := toolRegistry.NewToolRegistry()

	// mcpManagerServer.AddEnabledMCPServers(myTools)

	// now add common tools

	myTools.AddTool(aiGenTools.NewAiGenerateImageEditTool())
	myTools.AddTool(aiGenTools.NewAiImageDescriptionTool())

	// myTools.AddTool(commonTools.NewToolFileContents())
	// myTools.AddTool(aiGenTools.NewAiGenerateImageTool())
	// avatarTools.AddTool(aiGenTools.NewUiLaunchTool())
	// myTools.AddTool(apiTools.NewApiGetTool())

	m.ToolRegistry = myTools
	m.Tools = myTools.GetCommonTools()
	logrus.Debugf("Tools Added to conscious memory %s", m.Name)

}

func (m *Memory) DumpTools() {
	m.ToolRegistry.DumpTools(m.Tools)
}

func (m *Memory) DumpMemory(w io.Writer) {
	logrus.Debug("-----------------------------------------------------------------")
	logrus.Debugf("Name: %s, Model: %s, Temp: %f", m.Name, m.Model, m.Temperature)
	logrus.Debugf("Num Tools %d", len(m.Tools))
	// logrus.Debugf("System Prompt: %s", m.SystemPrompt)
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
		default:
			logrus.Debug("Non standard message found in memory *****")
		}
		if w != nil {
			// write to an io writer, so UI client or non standard client
			// can view the output
			str := fmt.Sprintf("Name: %s, content: %s\n", role, content)
			w.Write([]byte(str))
		}
		logrus.Debugf("Name: %s, content: %s", role, content)
	}
	logrus.Debug("--------------------------  END  ----------------------------------")

}

// AppendInMessage()
// Append IN Message, take the prompt and add to User Type message
// append to memory, only append user and assistant message to memory
func (m *Memory) AppendInMessage(in *pb.ContextMesssage) {
	// reset current context
	m.CurrentContext = &pb.SingleContext{
		In: in,
	}
	m.AppendMessage(openai.UserMessage(in.Prompt))
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
	m.AppendMessage(openai.AssistantMessage(out.Response))
}

// AppendMemory
// append to memory, only append user and assistant message to memory
func (m *Memory) AppendMessage(message openai.ChatCompletionMessageParamUnion) {

	// there is a change for user prompt to get duplicated
	// check if the current message is user message, then compare with last 2 messages
	if message.OfUser != nil {
		lastMessage := m.Messages[len(m.Messages)-1]
		if lastMessage.OfUser != nil {
			if message.OfUser.Content.OfString == lastMessage.OfUser.Content.OfString {
				logrus.Infof("user message duplication prevented %s", message.OfUser.Content.OfString)
				return
			}
		}
		if lastMessage.OfTool != nil && len(m.Messages) >= 2 {
			// if the last message was a tool message, lets compare the previous message
			lastMessage = m.Messages[len(m.Messages)-2]
			if lastMessage.OfUser != nil {
				if message.OfUser.Content.OfString == lastMessage.OfUser.Content.OfString {
					logrus.Infof("user message duplication prevented %s", message.OfUser.Content.OfString)
					return
				}
			}
		}
	}

	m.Messages = append(m.Messages, message)
	m.Save()
}

// // Load function to load from model file
// func (m *Memory) Load(name string) (*Memory, error) {

// }

func (m *Memory) Save() error {

	// scratch memory, which are unnamed memory
	// are not saved
	if m.Name == "" {
		return nil
	}
	fullFileName := dataLayer.GetFileAbsolutePath(MEMORY_DIRECTORY, m.Name)
	if err := dataLayer.SaveBinary(fullFileName, m); err != nil {
		logrus.Errorf("Cannot save file %s %v", fullFileName, err)
		return err
	}
	logrus.Infof("Saved memory to file %s", fullFileName)
	return nil
}

// LoadMemory
func LoadMemory(name string) (*Memory, error) {
	// Read profile back (Go automatically infers type T from your variable assignment)
	fullFileName := dataLayer.GetFileAbsolutePath(MEMORY_DIRECTORY, name)
	localMemory, err := dataLayer.ReadBinary[Memory](fullFileName)
	if err != nil {
		logrus.Errorf("error reading memory file %s, %v", name, err)
		return nil, err
	}

	// override model from load with settings value
	localMemory.Model = settings.Global.GetValue("LLM_MODEL")
	return &localMemory, nil
}
