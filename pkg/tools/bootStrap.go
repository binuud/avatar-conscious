package tools

import (
	"encoding/gob"

	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	"github.com/binuud/avatar-conscious/pkg/tools/apiTools"
	"github.com/binuud/avatar-conscious/pkg/tools/commonTools"
)

// var GlobalCommonToolsRegistry *toolRegistry.ToolRegistry
// var GlobalCommonTools []openai.ChatCompletionToolUnionParam

// func Bootstrap() {

// 	GlobalCommonToolsRegistry = nil
// 	GlobalCommonToolsRegistry = toolRegistry.NewToolRegistry()
// 	GlobalCommonToolsRegistry.AddTool(commonTools.NewToolFileContents())
// 	GlobalCommonToolsRegistry.AddTool(apiTools.NewApiGetTool())

// 	// GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageTool())
// 	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageEditTool())
// 	GlobalCommonToolsRegistry.AddTool(aiGenTools.NewAiGenerateImageDescriptionTool())
// 	// GlobalCommonToolsRegistry.AddTool(aiGenTools.NewUiLaunchTool())

// 	GlobalCommonTools = GlobalCommonToolsRegistry.GetCommonTools()

// }

func init() {
	gob.Register(commonTools.ToolFileContents{})
	gob.Register(apiTools.ApiGetTool{})
	gob.Register(aiGenTools.AiGenerateImageTool{})
	gob.Register(aiGenTools.AiGenerateImageEditTool{})
	gob.Register(aiGenTools.AiImageDescribeTool{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]string{})

}
