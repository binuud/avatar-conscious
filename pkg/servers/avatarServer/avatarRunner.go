package avatarServer

import (
	"fmt"
	"strings"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	pbConcsious "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// Run()
// iterates throught the steps, using llm
// updates the memory with user and assistant messages
// so tools down the chain understand the context
func (a *AvatarWrapper) Run() error {

	// iterate over the steps, as per instructctions on the avatar md file
	for i, step := range a.Avatar.Steps {
		_, err := a.RunStep(i, step)
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
	a.Memory.Tools = tools.GlobalCommonTools
	logrus.Infof("Running Avatar %s, with Model : %s", a.Meta.Name, a.Memory.Model)
	logrus.Infof("Description %s", a.Meta.Description)
}

// RunStep
// run a particular step, returns the outContext message
// update memory
func (a *AvatarWrapper) RunStep(i int, step *pb.Step) (*pbConcsious.ContextMesssage, error) {

	out := &pbConcsious.ContextMesssage{}
	logrus.Infof("%s (%d) : %s", step.Name, i+1, step.Prompt)
	a.Memory.AppendMemory(openai.UserMessage(fmt.Sprintf("Step No: %d, %s", i+1, step.Prompt)))

	a.Memory.DumpMemory(nil)

	// handle slash commands if any in the prompt
	if strings.HasPrefix(step.Prompt, "/") {
		err := a.RunSlashCommand(step.Prompt, out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}

	completions, err := llm.MakeLLMSingleToolRequest(a.Memory.Messages, a.Memory.Model,
		a.Memory.Tools, a.Memory.Temperature)
	if err != nil {
		logrus.Errorf("client: error llm chat request: %v \n", err)
		return nil, err
	}

	// check if any tool has to be called
	// if so call hanldle the tool and return the output response
	msg := completions.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		// a.Memory.AppendMemory(openai.ToolMessage(msg.ToolCalls[0].RawJSON(), msg.ToolCalls[0].ID))
		out, err = tools.GlobalCommonToolsRegistry.HandleToolCalls(msg.ToolCalls, step.Prompt)
		if err != nil {
			logrus.Errorf("HandleToolCalls: %v \n", err)
			return nil, err
		}
		logrus.Infof("Tool %s called, Out Msg : %s", msg.ToolCalls[0].Function.Name, out)
		// handle special case ui launch tool
		if msg.ToolCalls[0].Function.Name == "get_ui_launch_tool" {
			// if it is a show ui call
			// Action will be SHOW-UI
			// data will be id of the UI
			out.Action = "SHOW-UI"
			logrus.Infof("RunStep : SHOW-UI called with id %s", out.Data)
			a.Memory.AppendMemory(openai.ToolMessage(fmt.Sprintf("Avatar UI Interface with id %s shown to the user by the tool function", out.Data), msg.ToolCalls[0].ID))
			a.Memory.AppendMemory(openai.DeveloperMessage(fmt.Sprintf("Avatar UI Interface with id %s shown to the user by the tool function", out.Data)))
			// check if ui is present if not bail out
			_, err := a.CheckUiExists(out.Data)
			if err != nil {
				logrus.Errorf("RunStep %v", err)
				return nil, err
			}
		}
	} else {
		out.Response = msg.Content
		logrus.Infof("No Tools Called !!, Out Msg : %s", msg.Content)
	}
	return out, nil
}

func (a *AvatarWrapper) RunSlashCommand(prompt string, out *pbConcsious.ContextMesssage) error {
	// split into command and params
	result := strings.Fields(prompt)
	if len(result) == 2 {
		cmd := strings.ToLower(result[0])
		if cmd == "/launch-ui" {
			a.Memory.AppendMemory(openai.ToolMessage(fmt.Sprintf("Avatar UI Interface with id %s shown to the user by the tool function", out.Data), ""))
			a.Memory.AppendMemory(openai.DeveloperMessage(fmt.Sprintf("Avatar UI Interface with id %s shown to the user by the tool function", out.Data)))
			out.Action = "SHOW-UI"
			out.Data = result[1]
			// check if ui is present if not bail out
			_, err := a.CheckUiExists(out.Data)
			if err != nil {
				logrus.Errorf("RunStep %v", err)
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("cannot parse slash command")
}

// GetStepStatusesObj()
// creates an array of StepStatus, this is for running avatar
// and tracking the status
func (a *AvatarWrapper) GetStepStatusesObj() []*pb.StepStatus {
	numSteps := len(a.Avatar.Steps)
	stepStatuses := make([]*pb.StepStatus, 0)
	for i := 0; i < numSteps; i++ {
		stepStatuses = append(stepStatuses, &pb.StepStatus{Status: pb.AvatarStatus_STATUS_NOT_STARTED})
	}
	return stepStatuses
}

func (a *AvatarWrapper) CheckUiExists(id string) (int, error) {
	for i, ui := range a.UiCollection.Ui {
		if ui.Name == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("UI component with id %s not found", id)
}
