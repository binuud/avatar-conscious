package mcpHandler

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpCommandHandler struct {
	name  string
	Tools []*mcp.Tool `json:"tools"`
}

// sample "mcp-greet-server"
func NewMcpCommandHandler(mcpServerName string) (*McpCommandHandler, error) {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout.
	transport := &mcp.CommandTransport{Command: exec.Command(mcpServerName)}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	toolList, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}
	// for _, tool := range toolList.Tools {
	// 	logrus.Debugf("Tool Name %s \n", tool.Name)
	// 	logrus.Debugf("Tool Name %s \n", tool.Description)
	// 	logrus.Debugf("Tool Name %s \n", tool.InputSchema)
	// 	logrus.Debugf("Tool Name %s \n", tool.Meta)
	// }

	return &McpCommandHandler{
		name:  mcpServerName,
		Tools: toolList.Tools,
	}, nil
}

func HandleMcp(mcpServerName string, mcpFunction string, mcpParams any) (string, error) {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout.
	transport := &mcp.CommandTransport{Command: exec.Command(mcpServerName)}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return "", fmt.Errorf("tool failed %v", err)
	}
	defer session.Close()
	// Call a tool on the server.
	params := &mcp.CallToolParams{
		Name:      mcpFunction,
		Arguments: mcpParams,
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		return "", fmt.Errorf("CallTool failed: %v", err)
	}
	if res.IsError {
		return "", fmt.Errorf("tool failed %v", res)
	}
	out := ""
	for _, c := range res.Content {
		out = out + c.(*mcp.TextContent).Text
		log.Print(c.(*mcp.TextContent).Text)
	}

	return out, err
}
