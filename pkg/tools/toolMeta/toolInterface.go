package toolMeta

// Binu Udayakumar
// binu@dronasys.com

// All tools which will be registered to Avatar
// needs to implement this interface

import "github.com/openai/openai-go/v3"

type CommonToolInterface interface {

	// Name() returns the tool name
	// this is the name which is identified by LLM during tool call decision
	Name() string

	// Description()
	// returns the description
	// which is used by LLM, to decide if this tool has to be called,
	// for the given promp
	Description() string

	// InputSchema()
	// returns the openai function schema object for this tool's input object.
	InputSchema() string

	// InputSchemaAsOpenAiParam()
	// returns the openai function schema object for this tool's input object.
	InputSchemaAsOpenAiParam() openai.FunctionParameters

	// Handler(string)
	// makes the required call, bash, web fetch, curl etc,
	// this has to be implemented by the tool
	// please validated the input before making the tool call
	Handler(string) (string, error)
}
