package avatarServer

import (
	"fmt"
	"io"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarStoreServer"
	"github.com/sirupsen/logrus"
)

const (
	AVATAR_DIRECTORY       = "avatars"
	AVATAR_STORE_DIRECTORY = "avatarStore"
	SYSTEM_PROMPT          = `Review the provided instruction structure and extract a JSON list of all steps with prompts if needed,
														Steps should be prefixed in format {{Step [num]: [description]"}}
														Schema is provided in the context.
														do not parse any other properties not mentioned in the schema. Keep instructions in prompt field intact, 
														since llm will have to use the prompts as is, to complete the task later.
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
	return datalayer.GetFileAbsolutePath(AVATAR_STORE_DIRECTORY, a.Meta.Name)
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

	filePath := datalayer.GetFileAbsolutePath(AVATAR_DIRECTORY, name+".avatar")

	// Read AvatarWrapper back (Go automatically infers type T from your variable assignment)
	localAvatar, err := datalayer.ReadBinary[AvatarWrapper](filePath)
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

	filePath := datalayer.GetFileAbsolutePath(AVATAR_DIRECTORY, a.Avatar.Meta.Name+".avatar")
	// Save list of configs
	if err := datalayer.SaveBinary(filePath, a); err != nil {
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
