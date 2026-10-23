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
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	AVATAR_DIRECTORY       = "avatars"
	AVATAR_STORE_DIRECTORY = "avatarStore"
	SYSTEM_PROMPT          = `Review the provided instruction structure and extract a JSON list of all steps with prompts if needed,
														do not parse any other properties not mentioned in the schema`
)

type AvatarWrapper struct {
	*pb.Avatar
	FileName     string
	SystemPrompt string
	Memory       *memory.Memory
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

	logrus.Infof("Complete prompt : %s ", completePrompt)
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

	fmt.Println("\n output:")
	fmt.Println(output)

	stepsFromLLM := &pb.Steps{
		Steps: make([]*pb.Step, 0),
	}

	if err := json.Unmarshal([]byte(output), stepsFromLLM); err != nil {
		logrus.Fatalf("Avatar Build: Cannot unmarshall %v", err)
		return err
	}

	fmt.Printf("Struct %+v", stepsFromLLM)
	if stepsFromLLM != nil {
		a.Avatar.Steps = stepsFromLLM.Steps
	} else {
		logrus.Errorf("llm not able to build steps from %s.md", a.Meta.Name)
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
	fmt.Println(string(output))
	// Unmarshal the JSON bytes into your map
	err = json.Unmarshal(output, &schemaMap)
	if err != nil {
		fmt.Printf("Error unmarshalling schema: %v\n", err)
		return nil, err
	}
	return schemaMap, nil
}

func (a *AvatarWrapper) Run() {
	logrus.Infof(" >> Running Avatar %s \n", a.Meta.Name)
	logrus.Infof("		Description %s \n", a.Meta.Description)
	for _, item := range a.Avatar.Steps {
		logrus.Infof(" >>>> Running Step %s \n", item.Name)
		logrus.Infof("			Description %s \n", item.Description)
	}
}

func (a *AvatarWrapper) Print(w io.Writer) {
	if w == nil {
		return
	}
	str := fmt.Sprintf("Name: %s\nDescription %s\n\n", a.Meta.Name, a.Meta.Description)

	for _, item := range a.Avatar.Steps {
		str += fmt.Sprintf("%s Description : %s\nPrompt: %s\n\n", item.Name, item.Description, item.Prompt)
	}
	w.Write([]byte(str))
}
