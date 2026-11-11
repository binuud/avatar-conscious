package mcpHandler

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

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

	var session *mcp.ClientSession
	var err error
	url := strings.ToLower(mcpServerName)
	isHTTP := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
	if isHTTP {
		transport := &mcp.StreamableClientTransport{
			Endpoint: url,
		}
		// Connect to a server over stdin/stdout.
		session, err = client.Connect(ctx, transport, nil)
		if err != nil {
			return nil, err
		}
	} else {
		transport := &mcp.CommandTransport{Command: exec.Command(mcpServerName)}
		// Connect to a server over stdin/stdout.
		session, err = client.Connect(ctx, transport, nil)
		if err != nil {
			return nil, err
		}
	}

	defer session.Close()

	toolList, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &McpCommandHandler{
		name:  mcpServerName,
		Tools: toolList.Tools,
	}, nil
}

func HandleMcp(mcpServerName string, mcpFunction string, mcpParams any) (string, error) {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)
	var session *mcp.ClientSession
	var err error
	// Connect to a server over stdin/stdout.
	url := strings.ToLower(mcpServerName)
	isHTTP := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
	if isHTTP {
		transport := &mcp.StreamableClientTransport{
			Endpoint: url,
		}
		// Connect to a server over stdin/stdout.
		session, err = client.Connect(ctx, transport, nil)
		if err != nil {
			return "", err
		}
	} else {
		transport := &mcp.CommandTransport{Command: exec.Command(mcpServerName)}
		// Connect to a server over stdin/stdout.
		session, err = client.Connect(ctx, transport, nil)
		if err != nil {
			return "", err
		}
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
