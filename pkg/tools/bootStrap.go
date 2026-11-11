package tools

import (
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
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
	mcpName := "mcp-greet-server"
	mcp, err := mcpHandler.NewMcpCommandHandler(mcpName)
	if err != nil {
		logrus.Fatalf("Cannot get tools from mcp server/commnad %s : %v", mcpName, err)
	}

	// for each tool in mcp create an entry in GlobalTools
	for _, tool := range mcp.Tools {
		mcpTool := mcpCommandTool.NewMCPCommandTool(mcpName, tool)
		GlobalCommonToolsRegistry.AddMCPTool(mcpName, mcpTool)
	}

	GlobalCommonTools = GlobalCommonToolsRegistry.GetCommonTools()

}
