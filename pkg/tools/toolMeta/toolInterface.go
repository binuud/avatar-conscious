package toolMeta

import "github.com/openai/openai-go/v3"

type CommonToolInterface interface {

	// Name returns the predefind tool name (e.g. "Bash", "FileRead", "WebFetch").
	Name() string

	Description() string

	// InputSchema returns the JSON Schema for this tool's input object.
	InputSchema() string

	InputSchemaAsOpenAiParam() openai.FunctionParameters

	Handler(string) (string, error)

	// Return -1 for no limit.
	// MaxResultSizeChars() int
}
