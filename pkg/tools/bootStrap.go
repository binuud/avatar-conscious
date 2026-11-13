package tools

import (
	"github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/servers/mcpManagerServer"
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	apitools "github.com/binuud/avatar-conscious/pkg/tools/apiTools"
	"github.com/binuud/avatar-conscious/pkg/tools/commonTools"
	"github.com/binuud/avatar-conscious/pkg/tools/mcpCommandTool"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

var GlobalCommonToolsRegistry = toolRegistry.NewToolRegistry()
var GlobalCommonTools []openai.ChatCompletionToolUnionParam

func Bootstrap() {

	GlobalCommonToolsRegistry.AddTool(commonTools.NewToolFileContents())
	GlobalCommonToolsRegistry.AddTool(apitools.NewApiGetTool())
	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageTool())
	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageEditTool())
	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageDescriptionTool())
	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewUiLaunchTool())

	// now add mcp tools
	mcpServers, err := mcpManagerServer.Load()
	if err == nil {
		for _, mcpServer := range mcpServers.Configs {

			if mcpServer.Enabled == false {
				continue
			}

			// only add mcp server which are enabled
			if mcpServer.McpType == mcpManager.MCPTransportType_MCP_TYPE_STDIO {
				mcp, err := mcpHandler.NewMcpCommandHandler(mcpServer.Name, mcpServer.Command, "")
				if err != nil {
					logrus.Fatalf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
				}

				// for each tool in mcp create an entry in GlobalTools
				for _, tool := range mcp.Tools {
					mcpTool := mcpCommandTool.NewMCPCommandTool(mcpServer.Command, tool)
					GlobalCommonToolsRegistry.AddMCPTool(mcpServer.Name, mcpTool)
				}
			}
			if mcpServer.McpType == mcpManager.MCPTransportType_MCP_TYPE_HTTP {
				mcp, err := mcpHandler.NewMcpCommandHandler(mcpServer.Name, "", mcpServer.Url)
				if err != nil {
					logrus.Fatalf("Cannot get tools from mcp server/commnad %s : %v", mcpServer, err)
				}

				// for each tool in mcp create an entry in GlobalTools
				for _, tool := range mcp.Tools {
					mcpTool := mcpCommandTool.NewMCPCommandTool(mcpServer.Command, tool)
					GlobalCommonToolsRegistry.AddMCPTool(mcpServer.Name, mcpTool)
				}
			}
		}
	}
	// "http://localhost:1080/mockserver/mcp"

	GlobalCommonTools = GlobalCommonToolsRegistry.GetCommonTools()

}
