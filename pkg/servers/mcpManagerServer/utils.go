package mcpManagerServer

import (
	"github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"

	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/sirupsen/logrus"
)

// AddMcpServerToTool
// add mcp server name to the tool
// on handling, we search the enabled list of mcp servers for matching server name
// then run the tool on that server
func AddMcpServerToTool(mcpServer *mcpManager.MCPServerConfig, linkToolRegisty *toolRegistry.ToolRegistry) error {
	// only add mcp server which are enabled
	if mcpServer.Enabled == false {
		return nil
	}
	logrus.Debugf("AddMcpServerToTool - %s %+v", mcpServer.Name, mcpServer)

	if mcpServer.McpType == mcpManager.MCPTransportType_MCP_TYPE_STDIO {
		mcp, err := mcpHandler.NewMcpCommandHandler(mcpServer.Name, mcpServer.Command, "")
		if err != nil {
			logrus.Errorf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
			return err
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
			logrus.Errorf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
			return err
		}

		// for each tool in mcp create an entry in GlobalTools
		for _, tool := range mcp.Tools {
			mcpTool := NewMCPCommandTool(mcpServer.Name, tool)
			linkToolRegisty.AddMCPTool(mcpServer.Name, mcpTool)
		}
	}
	return nil
}

func AddEnabledMCPServers(linkToolRegisty *toolRegistry.ToolRegistry) {
	// now add mcp tools
	mcpServers := NewMcpServers()
	for _, mcpServer := range mcpServers.Servers.Configs {

		AddMcpServerToTool(mcpServer, linkToolRegisty)

	}
}
