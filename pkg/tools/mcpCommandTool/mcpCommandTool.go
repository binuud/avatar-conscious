package mcpCommandTool

// Binu Udayakumar
// binu@dronasys.com
// Do not change this file without my approval

// Tool - fileContents
// Get the contents of a file, whose path is given by LLM
// This can be used as reference to generate other others tools

import (
	"fmt"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"

	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/servers/mcpManagerServer"

	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

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

	// Safe type assertion to map[string]any
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

	mcpServers := mcpManagerServer.NewMcpServers()
	// since we do not have UUID in tooling, pass empty string for UUID
	outMsg, err := mcpServers.HandleTool("", t.name, t.mcpTool.Name, in.Prompt)
	if err != nil {
		return nil, err
	}

	return &pb.ContextMesssage{
		Uuid:     in.Uuid,
		Prompt:   in.Prompt,
		Response: outMsg,
	}, nil

}

// AddMcpServerToTool
// add mcp server name to the tool
// on handling, we search the enabled list of mcp servers for matching server name
// then run the tool on that server
func AddMcpServerToTool(mcpServer *mcpManager.MCPServerConfig, linkToolRegisty *toolRegistry.ToolRegistry) {
	// only add mcp server which are enabled
	if mcpServer.Enabled == false {
		return
	}
	logrus.Debugf("AddMcpServerToTool - %s %+v", mcpServer.Name, mcpServer)

	if mcpServer.McpType == mcpManager.MCPTransportType_MCP_TYPE_STDIO {
		mcp, err := mcpHandler.NewMcpCommandHandler(mcpServer.Name, mcpServer.Command, "")
		if err != nil {
			logrus.Fatalf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
		}

		// for each tool in mcp create an entry in GlobalTools
		for _, tool := range mcp.Tools {
			mcpTool := NewMCPCommandTool(mcpServer.Name, tool)
			linkToolRegisty.AddMCPTool(mcpServer.Name, mcpTool)
		}
	}
	if mcpServer.McpType == mcpManager.MCPTransportType_MCP_TYPE_HTTP {
		logrus.Debugf("AddMcpServerTool remote tool %s", mcpServer.Name)
		mcp, err := mcpHandler.NewMcpCommandHandler(mcpServer.Name, "", mcpServer.Url)
		if err != nil {
			logrus.Fatalf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
		}

		// for each tool in mcp create an entry in GlobalTools
		for _, tool := range mcp.Tools {
			mcpTool := NewMCPCommandTool(mcpServer.Name, tool)
			linkToolRegisty.AddMCPTool(mcpServer.Name, mcpTool)
		}
	}

}
