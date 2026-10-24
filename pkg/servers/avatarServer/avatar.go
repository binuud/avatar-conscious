package avatarServer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarStoreServer"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	AVATAR_DIRECTORY       = "avatars"
	AVATAR_STORE_DIRECTORY = "avatarStore"
	SYSTEM_PROMPT          = `Review the provided instruction structure and extract a JSON list of all steps with prompts if needed,
														do not parse any other properties not mentioned in the schema. Keep instructions in prompt field intact, 
														since llm will have to use the prompts as is, to complete the task later.`
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

func (a *AvatarWrapper) GetStepsFromLLM() error {
	return a.Build()
}

func (a *AvatarWrapper) Load(name string) error {

	filePath := datalayer.GetFileAbsolutePath(AVATAR_DIRECTORY, name+".avatar")

	// Read AvatarWrapper back (Go automatically infers type T from your variable assignment)
	localAvatar, err := datalayer.ReadBinary[AvatarWrapper](filePath)
	if err != nil {
		return err
	}
	a.Avatar.Steps = localAvatar.Steps
	return nil
}

func (a *AvatarWrapper) Save() error {

	if a.Avatar.Meta.Name == "" {
		return fmt.Errorf("avatar name missing, cannot save")
	}

	filePath := datalayer.GetFileAbsolutePath(AVATAR_DIRECTORY, a.Avatar.Meta.Name+".avatar")
	// Save list of configs
	if err := datalayer.SaveBinary(filePath, a); err != nil {
		logrus.Errorf("Cannot save avatar file %s %v", filePath, err)
		return err
	}
	logrus.Infof("Saved avatar to file %s", filePath)
	return nil
}

// Build()
// builds the avatar file, creates the necessary steps
// to run the avatar using llm and the avatar with the steps
func (a *AvatarWrapper) Build() error {

	completePrompt := fmt.Sprintf(
		"Here is the avatar context file (%s):\n\n```markdown\n%s\n```\n\n",
		a.FileName,
		a.Content,
	)

	a.Memory.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(a.SystemPrompt),
		openai.UserMessage(completePrompt),
	}

	stepsSchema, err := a.getAvatarStepsJsonSchemaMap()
	if err != nil {
		logrus.Errorf("parseAvatar: %v", err)
		return err
	}

	// logrus.Infof("Complete prompt : %s ", completePrompt)
	// logrus.Infof("StepsSchema : %v", stepsSchema)

	// now use LLM to get the steps, that this avatar needs to run
	output, err := llm.MakeLLMSchemaResponseRequest(completePrompt,
		a.SystemPrompt,
		a.Memory.Model,
		stepsSchema)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return err
	}

	// fmt.Println("\n output:")
	// fmt.Println(output)

	stepsFromLLM := &pb.Steps{
		Steps: make([]*pb.Step, 0),
	}

	if err := json.Unmarshal([]byte(output), stepsFromLLM); err != nil {
		logrus.Fatalf("Avatar Build: Cannot unmarshall %v", err)
		return err
	}

	fmt.Printf("Struct %+v", stepsFromLLM)
	if len(stepsFromLLM.Steps) == 0 {
		logrus.Errorf("llm not able to build steps from %s.md", a.Meta.Name)
	} else {
		a.Avatar.Steps = stepsFromLLM.Steps
	}

	return nil
}

// getAvatarStepsJsonSchemaMap()
// converts the Avatar Struct to map[string] any
// this is needed, because openAiAPI expects struct schema as a map object
func (a *AvatarWrapper) getAvatarStepsJsonSchemaMap() (map[string]any, error) {

	var schemaMap map[string]any

	schemaJSON, err := jsonschema.For[pb.Steps](nil)
	if err != nil {
		log.Fatalf("Failed to run  jsonschema: %v", err)
	}

	output, err := json.MarshalIndent(schemaJSON, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal schema: %v", err)
	}
	// fmt.Println(string(output))
	// Unmarshal the JSON bytes into your map
	err = json.Unmarshal(output, &schemaMap)
	if err != nil {
		fmt.Printf("Error unmarshalling schema: %v\n", err)
		return nil, err
	}
	return schemaMap, nil
}

// Run()
// iterates throught the steps, using llm
// updates the memory with user and assistant messages
// so tools down the chain understand the context
func (a *AvatarWrapper) Run() error {

	// iterate over the steps, as per instructctions on the avatar md file
	for i, step := range a.Avatar.Steps {
		err := a.RunStep(i, step)
		if err != nil {
			return err
		}
	}
	return nil
}

// Start
// start the avatar, creates the memory
// waits for run by step command
func (a *AvatarWrapper) Start() {
	// create an use and throw memory
	a.Memory = memory.NewMemory("")
	a.Memory.Tools = tools.GlobalCommonTools.GetCommonTools()
	logrus.Infof("Running Avatar %s, with Model : %s", a.Meta.Name, a.Memory.Model)
	logrus.Infof("Description %s", a.Meta.Description)
}

// RunStep
// run a particular step
func (a *AvatarWrapper) RunStep(i int, step *pb.Step) error {
	logrus.Infof("%s (%d) : %s", step.Name, i+1, step.Prompt)
	a.Memory.Messages = append(a.Memory.Messages, openai.UserMessage(step.Prompt))

	completions, err := llm.MakeLLMSingleToolRequest(a.Memory.Messages, a.Memory.Model,
		a.Memory.Tools, a.Memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return err
	}

	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		out, err := tools.GlobalCommonTools.HandleToolCalls(msg.ToolCalls)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return err
		}
		logrus.Infof("Out Msg : %s", out)
		a.Memory.Messages = append(a.Memory.Messages, openai.AssistantMessage(out.Response))
	} else {
		logrus.Infof("Out Msg : %s", msg.Content)
	}
	return nil
}

// GetStepStatusesObj()
// creates an array of StepStatus, this is for running avatar
// and tracking the status
func (a *AvatarWrapper) GetStepStatusesObj() []*pb.StepStatus {
	numSteps := len(a.Avatar.Steps)
	stepStatuses := make([]*pb.StepStatus, 0)
	for i := 0; i < numSteps; i++ {
		stepStatuses = append(stepStatuses, &pb.StepStatus{Status: pb.StepStatus_NotStarted})
	}
	return stepStatuses
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
