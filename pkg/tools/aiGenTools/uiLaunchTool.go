package aiGenTools

import (
	"encoding/json"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// LLM returns a json string
// which has to be parsed into a equivalent structure
// this provides the input structure
type UiLaunchTool_Input struct {
	// refined text prompt
	Text string `json:"text"`
}

// implements the CommonToolInterface
type UiLaunchTool struct {
	toolMeta.CommonToolInterface
}

func NewUiLaunchTool() *UiLaunchTool {
	return &UiLaunchTool{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t UiLaunchTool) Name() string {
	return "get_ui_launch_tool"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t UiLaunchTool) Description() string {
	return "Launch an Avatar UI interface."
}

func (t UiLaunchTool) InputSchema() string {
	return "Launch an Avatar UI interface."
}

func (t UiLaunchTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"text": map[string]string{
				"type":        "string",
				"description": "Launch an Avatar UI interface. Instructions from the prompt is saved here.",
			},
		},
		"required": []string{"text"},
	}
}

func (t UiLaunchTool) Handler(in *pb.ContextMesssage) (*pb.ContextMesssage, error) {
	logrus.Infof("UiLaunchTool Handler Arguments received %v", in.Prompt)
	var inpObj UiLaunchTool_Input
	if err := json.Unmarshal([]byte(in.Prompt), &inpObj); err != nil {
		return nil, err
	}
	logrus.Infof("Arguments parsed into UiLaunchTool_Input structure %v", in)

	return &pb.ContextMesssage{
		Uuid:   in.Uuid,
		Prompt: in.Prompt,
		Action: inpObj.Text,
	}, nil
}
