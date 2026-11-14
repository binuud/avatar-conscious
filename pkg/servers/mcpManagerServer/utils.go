package mcpManagerServer

import (
	"github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	pbMcp "github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/tools/mcpCommandTool"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/sirupsen/logrus"
)

// AddMcpServerToTool
// add mcp server to tools, only if it is enabled
func AddMcpServerToTool(mcpServer *pbMcp.MCPServerConfig, linkToolRegisty *toolRegistry.ToolRegistry) {
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
			mcpTool := mcpCommandTool.NewMCPCommandTool(mcpServer.Command, tool)
			linkToolRegisty.AddMCPTool(mcpServer.Name, mcpTool)
		}
	}
	if mcpServer.McpType == mcpManager.MCPTransportType_MCP_TYPE_HTTP {
		logrus.Debugf("AddMcpServerToTool remote tool %s", mcpServer.Name)
		mcp, err := mcpHandler.NewMcpCommandHandler(mcpServer.Name, "", mcpServer.Url)
		if err != nil {
			logrus.Fatalf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
		}

		// for each tool in mcp create an entry in GlobalTools
		for _, tool := range mcp.Tools {
			mcpTool := mcpCommandTool.NewMCPCommandTool(mcpServer.Command, tool)
			linkToolRegisty.AddMCPTool(mcpServer.Name, mcpTool)
		}
	}

}
