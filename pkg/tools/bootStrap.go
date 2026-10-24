package tools

import (
	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	apitools "github.com/binuud/avatar-conscious/pkg/tools/apiTools"
	"github.com/binuud/avatar-conscious/pkg/tools/commonTools"
	"github.com/binuud/avatar-conscious/pkg/tools/toolRegistry"
)

var GlobalCommonTools = toolRegistry.NewToolRegistry()

func Bootstrap() {

	GlobalCommonTools.AddTool(commonTools.NewToolFileContents())
	GlobalCommonTools.AddTool(apitools.NewApiGetTool())
	GlobalCommonTools.AddTool(aiGenTools.NewAiGenerateImageTool())
	GlobalCommonTools.AddTool(aiGenTools.NewAiGenerateImageDescriptionTool())

}
