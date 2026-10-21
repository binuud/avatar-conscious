package avatar

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	AVATAR_DIRECTORY       = "avatars"
	AVATAR_STORE_DIRECTORY = "avatarStore"
)

type Step struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

type Avatar struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
}

type AvatarWrapper struct {
	Avatar
	FileName     string
	SystemPrompt string
	Memory       *memory.Memory
}

func NewAvatar(name string) *AvatarWrapper {

	avatar := &AvatarWrapper{
		FileName:     name,
		Memory:       memory.NewMemory("memory-" + name),
		SystemPrompt: "Review the provided instruction structure and extract a JSON list of all steps with prompts if needed, api endpoints and commands to be run.",
	}

	err := avatar.Load(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, check if the avatar was built, %v", err)
		return nil
	}
	return avatar
}

func NewAvatarFromStore(fileName string) *AvatarWrapper {

	avatar := &AvatarWrapper{
		FileName:     fileName,
		Memory:       memory.NewMemory("memory-" + fileName),
		SystemPrompt: "Review the provided instruction structure and extract a JSON list of all steps with prompts if needed, api endpoints and commands to be run.",
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
	*a = localAvatar
	return nil
}

func (a *AvatarWrapper) Save() error {

	if a.Avatar.Name == "" {
		return fmt.Errorf("avatar name missing, cannot save")
	}

	filePath := datalayer.GetFileAbsolutePath(AVATAR_DIRECTORY, a.Avatar.Name+".avatar")
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

	// read the file
	filePath := datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, a.FileName)
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("❌ Error reading %s: %v\n", a.FileName, err)
		return err
	}

	avatarContent := string(fileBytes)
	completePrompt := fmt.Sprintf(
		"Here is the avatar context file (%s):\n\n```markdown\n%s\n```\n\n",
		a.FileName,
		avatarContent,
	)

	a.Memory.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(a.SystemPrompt),
		openai.UserMessage(completePrompt),
	}

	schema, err := a.getAvatarJsonSchemaMap()
	if err != nil {
		logrus.Errorf("parseAvatar: %v", err)
		return err
	}

	logrus.Infof("Complete prompt : %s ", completePrompt)
	logrus.Infof("Schema : %s ", schema)

	output, err := llm.MakeLLMSchemaResponseRequest(completePrompt,
		a.SystemPrompt,
		a.Memory.Model,
		schema)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return err
	}

	var in Avatar
	if err := json.Unmarshal([]byte(output), &in); err != nil {
		logrus.Fatalf("Cannot unmarshall %v", err)
		return err
	}
	fmt.Println("\n✨ output:")
	fmt.Println(output)
	fmt.Printf("Struct %+v", in)
	a.Avatar = in
	return nil
}

// getAvatarJsonSchemaMap()
// converts the Avatar Struct to map[string] any
// this is needed, because openAiAPI expects struct schema as a map object
func (a *AvatarWrapper) getAvatarJsonSchemaMap() (map[string]any, error) {

	var schemaMap map[string]any

	schemaJSON, err := jsonschema.For[Avatar](nil)
	if err != nil {
		log.Fatalf("Failed to run  jsonschema: %v", err)
	}

	output, err := json.MarshalIndent(schemaJSON, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal schema: %v", err)
	}

	// Unmarshal the JSON bytes into your map
	err = json.Unmarshal(output, &schemaMap)
	if err != nil {
		fmt.Printf("Error unmarshalling schema: %v\n", err)
		return nil, err
	}
	return schemaMap, nil
}

func (a *AvatarWrapper) Run() {
	logrus.Infof(" >> Running Avatar %s \n", a.Name)
	logrus.Infof("		Description %s \n", a.Description)
	for _, item := range a.Avatar.Steps {
		logrus.Infof(" >>>> Running Step %s \n", item.Name)
		logrus.Infof("			Description %s \n", item.Description)
	}
}

func (a *AvatarWrapper) Print(w io.Writer) {
	if w == nil {
		return
	}
	str := fmt.Sprintf("Name: %s\nDescription %s\n\n", a.Name, a.Description)

	for _, item := range a.Avatar.Steps {
		str += fmt.Sprintf("Step Name: %s\nStep Description %s\n\n", item.Name, item.Description)
	}
	w.Write([]byte(str))
}
