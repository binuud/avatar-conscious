package mcpCommandTool

// Binu Udayakumar
// binu@dronasys.com
// Do not change this file without my approval

// Tool - fileContents
// Get the contents of a file, whose path is given by LLM
// This can be used as reference to generate other others tools

import (
	"encoding/json"
	"fmt"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"

	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// LLM returns a json string
// which has to be parsed into a equivalent structure
// this provides the input structure
type Tool_MCP_Input map[string]any

// path is the absolute path to read (required).

// implements the CommonToolInterface
type MCPCommandTool struct {
	name    string
	mcpTool *mcp.Tool
	toolMeta.CommonToolInterface
}

func NewMCPCommandTool(name string, tool *mcp.Tool) *MCPCommandTool {

	return &MCPCommandTool{
		name:    name,
		mcpTool: tool,
	}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t *MCPCommandTool) Name() string {
	return t.mcpTool.Name
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t *MCPCommandTool) Description() string {
	return t.mcpTool.Description
}

func (t *MCPCommandTool) InputSchema() string {
	return fmt.Sprint(t.mcpTool.InputSchema)
}

func (t *MCPCommandTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	// var schema = t.InputSchema()
	// var resultMap map[string]any
	// if err := json.Unmarshal([]byte(schema), &resultMap); err != nil {
	// 	logrus.Errorf("Cannot parse input schema %v", t.InputSchema())
	// 	logrus.Panicf("MCP convert InputSchemaAsOpenAiParam %s:%s, %v  ", t.name, t.mcpTool.Name, err)
	// }

	// 1. Safe type assertion to map[string]any
	resultMap, ok := t.mcpTool.InputSchema.(map[string]any)
	if !ok {
		logrus.Errorf("Cannot parse input schema %v", t.InputSchema())
		return nil
	}
	logrus.Infof("Parsed mcp schema %s: %v", t.name, resultMap)
	return openai.FunctionParameters(resultMap)
}

func (t *MCPCommandTool) Handler(in *pb.ContextMesssage) (*pb.ContextMesssage, error) {

	logrus.Infof("MCPCommandTool handler Arguments received %v", in.Prompt)
	var inObj Tool_MCP_Input
	if err := json.Unmarshal([]byte(in.Prompt), &inObj); err != nil {
		return nil, err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)
	logrus.Infof("Current tool %s %s", t.name, t.mcpTool.Name)
	outMsg, err := mcpHandler.HandleMcp(t.name, t.mcpTool.Name, inObj)
	if err != nil {
		return nil, err
	}

	return &pb.ContextMesssage{
		Uuid:     in.Uuid,
		Prompt:   in.Prompt,
		Response: outMsg,
	}, nil

}
