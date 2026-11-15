package tools

import (
	"github.com/binuud/avatar-conscious/pkg/servers/mcpManagerServer"
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	"github.com/binuud/avatar-conscious/pkg/tools/apiTools"
	"github.com/binuud/avatar-conscious/pkg/tools/commonTools"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
	"github.com/openai/openai-go/v3"
)

var GlobalCommonToolsRegistry = toolRegistry.NewToolRegistry()
var GlobalCommonTools []openai.ChatCompletionToolUnionParam

func Bootstrap() {

	GlobalCommonToolsRegistry.AddTool(commonTools.NewToolFileContents())
	GlobalCommonToolsRegistry.AddTool(apiTools.NewApiGetTool())

	// GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageTool())
	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageEditTool())
	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageDescriptionTool())
	// GlobalCommonToolsRegistry.AddTool(aiGenTools.NewUiLaunchTool())

	// now add mcp tools
	mcpServers := mcpManagerServer.NewMcpServers()
	for _, mcpServer := range mcpServers.Servers.Configs {

		mcpManagerServer.AddMcpServerToTool(mcpServer, GlobalCommonToolsRegistry)

	}

	GlobalCommonTools = GlobalCommonToolsRegistry.GetCommonTools()

}
