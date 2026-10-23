package avatarServer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	pbQuantumUI "github.com/binuud/avatar-conscious/gen/go/v1/quantumUI"
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
)

// type Step struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Prompt      string `json:"prompt"`
// }

// type Avatar struct {
// 	Uuid        string `json:"uuid"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Steps       []Step `json:"steps"`
// }

type AvatarWrapper struct {
	*pb.Avatar
	FileName     string
	SystemPrompt string
	Memory       *memory.Memory
}

func NewQuantumUi() *pbQuantumUI.QuantumUI {
	return &pbQuantumUI.QuantumUI{
		Name:        "First Pane",
		Description: "This is a test pane",
		Title:       "Avatar Pane",
		Components: []*pbQuantumUI.UiComponent{
			{
				Id:          "firstName",
				Label:       "First Name",
				PlaceHolder: "Enter your first name",
				Value:       "",
				Width:       50,
				Type:        pbQuantumUI.UiComponent_INPUT,
			},
			{
				Id:          "lastName",
				Label:       "Last Name",
				PlaceHolder: "Enter your last name",
				Value:       "",
				Width:       50,
				Type:        pbQuantumUI.UiComponent_INPUT,
			},
			{
				Value: "lorem ipsum dolor sit amet consectetur adipiscing elit optio ducimus temporibus ullamco ipsum dolor omnis aut et fugiat fugiat et voluptas et vero est autem minim possimus et libero facere imperdiet nihil consequatur eiusmod nostrud voluptas vero aliquip possimus duis qui dolorum et cupiditate elit minus quas ea voluptas consequatur",
				Width: 100,
				Type:  pbQuantumUI.UiComponent_PLAINTEXT,
			},
		},
	}
}

func NewAvatar(name string) *AvatarWrapper {

	avatar := &AvatarWrapper{
		FileName:     name,
		Memory:       memory.NewMemory("memory-" + name),
		SystemPrompt: "Review the provided instruction structure and extract a JSON list of all steps with prompts if needed, api endpoints and commands to be run.",
	}

	a, err := avatarStoreServer.ReadAvatarFromStore(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, check if the avatar was built, %v", err)
		return nil
	}

	// hard code UI for testing

	avatar.Avatar = a
	avatar.Ui = NewQuantumUi()

	// try reading the built avatar
	err = avatar.Load(name)
	if err != nil {
		logrus.Errorf("cannot load avatar, check if the avatar was built, %v", err)
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

	var in *pb.Avatar
	if err := json.Unmarshal([]byte(output), in); err != nil {
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

	schemaJSON, err := jsonschema.For[pb.Avatar](nil)
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
		str += fmt.Sprintf("Step Name: %s\nStep Description %s\n\n", item.Name, item.Description)
	}
	w.Write([]byte(str))
}
