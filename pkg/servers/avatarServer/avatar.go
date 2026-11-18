package avatarServer

import (
	"fmt"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarStoreServer"
	"github.com/sirupsen/logrus"
)

const (
	AVATAR_DIRECTORY       = "avatars"
	AVATAR_STORE_DIRECTORY = "avatarStore"
	SYSTEM_PROMPT          = `<role>
You are Avatar, an intelligent and careful generic assistant. 
You have access to external tools and resources provided via the Model Context Protocol (MCP). 
Your objective is to help the user accurately, safely, and efficiently by reasoning through their requests and executing the right tool sequence.
</role>

<core_principles>
1. **Grounding:** Never guess, fabricate, or hallucinate tool parameters. If a required argument is missing from the user's prompt, ask a direct clarifying question before calling the tool.
2. **Minimalism:** Only call tools when necessary. If a question can be answered from your existing knowledge or prior context without a tool, do not invoke one.
3. **Sequential Reasoning:** Think step-by-step. Before triggering any tool call, evaluate whether you have all required context and if this specific tool is the optimal choice for the current step.
4. **Error Recovery:** If a tool returns an error or empty result, analyze the failure message, adjust your approach or parameters, and try a different strategy instead of repeating the exact same failing call.
5. **No Speculation:** Never invoke a tool "just to see what it does." Every tool call must have a clear, goal-driven intent tied to the user's request.
</core_principles>

<workflow>
When a user gives a request:
1. **Analyze:** Parse the user's intent and identify what data or action is needed.
2. **Plan:** Determine if a tool is required. If multiple tools are needed, plan the execution order (e.g., fetch data first, then process or update).
3. **Execute:** Invoke the tool with strictly validated parameters matching the schema.
4. **Respond:** Once the tool output returns, synthesize a clear, concise, and helpful final answer for the user grounded strictly in the tool's response.
</workflow>

<safety_guardrails>
- Never leak raw API keys, internal system tokens, or sensitive credentials found in tool responses.
- If a user asks you to perform a destructive action (e.g., deleting a database record, dropping a table, or sending external communications), ask for explicit confirmation and state what parameters you intend to use *before* executing the tool.
</safety_guardrails>
`
)

type AvatarWrapper struct {
	*pb.Avatar
	FileName          string
	SystemPrompt      string
	Memory            *memory.Memory
	CurrentStep       *pb.Step
	CurrentStepNumber int
}

func NewAvatar(name string) *AvatarWrapper {

	avatar := &AvatarWrapper{
		FileName:     name,
		Memory:       memory.NewMemory("memory-" + name),
		SystemPrompt: SYSTEM_PROMPT,
	}

	a, err := avatarStoreServer.ReadAvatarFromStore(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, %v", err)
		return nil
	}

	// hard code UI for testing
	avatar.Avatar = a
	// avatar.Ui = NewQuantumUi()
	// try reading the built avatar
	err = avatar.Load(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, check if the avatar was built, %v", err)
	}

	return avatar

}

func NewAvatarFromStore(name string) *AvatarWrapper {

	avatar := &AvatarWrapper{
		FileName:     name,
		Memory:       memory.NewMemory("memory-" + name),
		SystemPrompt: SYSTEM_PROMPT,
	}

	a, err := avatarStoreServer.ReadAvatarFromStore(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, %v", err)
		return nil
	}

	avatar.Avatar = a
	err = avatar.Load(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, check if the avatar was built, %v", err)
	}

	return avatar
}

func (a *AvatarWrapper) GetAvatarDir() string {
	return dataLayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, a.Meta.Name)
}

// Get
// read the avatar file details from yaml files
func (a *AvatarWrapper) Get(name string) error {

	err := a.ReadAvatarUiFromYaml()
	if err != nil {
		logrus.Errorf("cannot parse or read ui collection yaml, %v", err)
	}
	err = a.ReadAvatarStepsFromYaml()
	if err != nil {
		logrus.Errorf("cannot parse or read steps yaml, %v", err)
	}
	return nil
}

// Load
// load avatar details from cached binary file
func (a *AvatarWrapper) Load(name string) error {

	filePath := dataLayer.GetFileAbsolutePath(AVATAR_DIRECTORY, name+".avatar")

	// Read AvatarWrapper back (Go automatically infers type T from your variable assignment)
	localAvatar, err := dataLayer.ReadBinary[pb.Avatar](filePath)
	if err != nil {
		return err
	}
	a.Avatar.UiCollection = localAvatar.UiCollection
	a.Avatar.Steps = localAvatar.Steps
	return nil
}

func (a *AvatarWrapper) Save() error {

	if a.Avatar.Meta.Name == "" {
		return fmt.Errorf("avatar name missing, cannot save")
	}

	filePath := dataLayer.GetFileAbsolutePath(AVATAR_DIRECTORY, a.Avatar.Meta.Name+".avatar")
	// Save list of configs
	if err := dataLayer.SaveBinary(filePath, a.Avatar); err != nil {
		logrus.Errorf("Cannot save avatar file %s %v", filePath, err)
		return err
	}
	logrus.Infof("Saved avatar to file %s", filePath)
	return nil
}

func (a *AvatarWrapper) Print(w io.Writer) {
	if w == nil {
		return
	}
	str := fmt.Sprintf("Name: %s\nDescription %s\n\n", a.Meta.Name, a.Meta.Description)

	for _, item := range a.Avatar.Steps {
		str = str + fmt.Sprintf("%s Description : %s\nPrompt: %s\n\n", item.Name, item.Description, item.Prompt)
	}
	w.Write([]byte(str))
}
