package avatarServer

import (
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// Build()
// builds the avatar file, creates the necessary steps
// to run the avatar using llm and the avatar with the steps
func (a *AvatarWrapper) Build() error {

	logrus.Infof("Building avatar %s", a.FileName)

	// unset previous build steps
	a.Steps = nil
	a.UiCollection = nil

	// try reading from yaml
	a.ReadAvatarStepsFromYaml()
	a.ReadAvatarUiFromYaml()

	if a.Steps == nil {
		// we do not have steps from the yaml file
		// try using LLM to get steps from the MD file
		err := a.BuildUsingLLM()
		if err != nil {
			return err
		}
	} else {
		logrus.Infof("User provided %s-steps.yaml file found for avatar, not building with llm ", a.Meta.Name)
	}

	a.Save()
	return nil
}

func (a *AvatarWrapper) BuildUsingLLM() error {
	userPrompt := fmt.Sprintf(
		"Here is the avatar context file (%s):\n\n```markdown\n%s\n```\n\n",
		a.FileName,
		a.Content,
	)

	buildMemory := memory.NewMemory()

	buildMemory.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(a.SystemPrompt),
		openai.UserMessage(userPrompt),
	}

	stepsSchema, err := a.getAvatarStepsJsonSchemaMap()
	if err != nil {
		logrus.Errorf("parseAvatar: %v", err)
		return err
	}

	// logrus.Infof("Complete prompt : %s ", completePrompt)
	// logrus.Infof("StepsSchema : %v", stepsSchema)

	// now use LLM to get the steps, that this avatar needs to run
	output, err := llm.MakeLLMSchemaResponseRequest(userPrompt,
		a.SystemPrompt,
		buildMemory.Model,
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

	if len(stepsFromLLM.Steps) == 0 {
		logrus.Errorf("llm not able to build steps from %s.md", a.Meta.Name)
	} else {
		a.Avatar.Steps = stepsFromLLM.Steps
	}
	logrus.Infof("Steps for avatar %s, parsed from MD file %+v", a.Meta.Name, stepsFromLLM)
	return nil

}

func (a *AvatarWrapper) ReadAvatarStepsFromYaml() error {

	avatarStepsFile := fmt.Sprintf("%s/%s-steps.yaml", a.GetAvatarDir(), a.Meta.Name)
	if utils.FileExists(avatarStepsFile) {
		steps := &pb.Steps{}
		err := datalayer.ReadProtoYAML(avatarStepsFile, steps)
		if err != nil {
			return err
		}
		a.Steps = steps.Steps
	}
	return nil
}

func (a *AvatarWrapper) ReadAvatarUiFromYaml() error {

	avatarUiFile := fmt.Sprintf("%s/%s-ui.yaml", a.GetAvatarDir(), a.Meta.Name)
	if utils.FileExists(avatarUiFile) {
		uiCollection := &pb.UiCollection{}
		err := datalayer.ReadProtoYAML(avatarUiFile, uiCollection)
		if err != nil {
			return err
		}
		a.UiCollection = uiCollection
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
