package toolRegistry

// Binu Udayakumar
// binu@dronasys.com
// Tool Registry
// Do Not Change, without my approval

// Tool Registry

// Toools have to be registered with Tool Registry
// all tool calls from llm will be handled by Tool Registry
// it can be shell commands, web fetch, web form submit
// everything has to go through Tool Regitry

import (
	"fmt"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// Tool Registry is a collection of tools
// Tools have to be defined in go
// and registered with tool registry
// Tools have to implement the toolMeta.CommonToolInterface interface
type ToolRegistry struct {
	Tools map[string]toolMeta.CommonToolInterface
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		Tools: make(map[string]toolMeta.CommonToolInterface),
	}
}

// AddTool()
// Tools have to be defined and added to registry
// Even skills can only call tools already registered with tool registry
func (t *ToolRegistry) AddTool(tool toolMeta.CommonToolInterface) {
	logrus.Debugf("AddedTool %s to ToolRegistry", tool.Name())
	t.Tools[tool.Name()] = tool
}

// AddTool()
// Tools have to be defined and added to registry
// Even skills can only call tools already registered with tool registry
func (t *ToolRegistry) GetTool(name string) (toolMeta.CommonToolInterface, error) {
	if tool, ok := t.Tools[name]; ok {
		return tool, nil
	}
	return nil, fmt.Errorf("tool with name %s, not registered", name)

}

// GetCommonTools()
// Builds the Tools Array for openAI api, from the registered set of tools
func (t *ToolRegistry) GetCommonTools() []openai.ChatCompletionToolUnionParam {

	myTools := make([]openai.ChatCompletionToolUnionParam, 0)

	for name, tool := range t.Tools {
		logrus.Debugf("Build tool - %s, %s", name, tool.Name())
		logrus.Debugf("Description - %s, %s", name, tool.Description())

		newToolFunction := openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        tool.Name(),
			Description: openai.String(tool.Description()),
			Parameters:  tool.InputSchemaAsOpenAiParam(),
		})
		myTools = append(myTools, newToolFunction)
	}
	return myTools

}

// HandleToolCalls()
// When LLM returns a set of tools to be handled,
// the handler runs each of the tool calls, with the arguments given by LLM
func (t *ToolRegistry) HandleToolCalls(toolCalls []openai.ChatCompletionMessageToolCallUnion) (*pb.ContextMesssage, error) {
	if len(toolCalls) < 1 {
		return nil, fmt.Errorf("No tools to be handled")
	}

	logrus.Infof("LLM has returned tool calls. Handling them...")

	// handle first tool call for now
	toolCall := toolCalls[0]
	logrus.Infof("Function %+v", toolCall.Function.Name)
	logrus.Infof("Args %+v", toolCall.Function.Arguments)
	myTool, err := t.GetTool(toolCall.Function.Name)
	if err != nil {
		return nil, err
	}
	in := &pb.ContextMesssage{
		Uuid:   utils.GenerateUUID(),
		Prompt: toolCall.Function.Arguments,
	}
	return myTool.Handler(in)
}
