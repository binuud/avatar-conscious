package mcpHandler

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
)

type McpCommandHandler struct {
	name  string
	Tools []*mcp.Tool `json:"tools"`
}

// sample "mcp-greet-server"
// NewMcpCommandHandler
// name is human readable snake case
// cmd or url can be provided, both should not be emtpy
func NewMcpCommandHandler(name string, cmd string, url string) (*McpCommandHandler, error) {

	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	var session *mcp.ClientSession
	var err error

	if cmd != "" {
		transport := &mcp.CommandTransport{Command: exec.Command(cmd)}
		// Connect to a server over stdin/stdout.
		session, err = client.Connect(ctx, transport, nil)
		if err != nil {
			return nil, err
		}
	} else if url != "" {
		lowerCaseUrl := strings.ToLower(url)
		isHTTP := strings.HasPrefix(lowerCaseUrl, "http://") || strings.HasPrefix(lowerCaseUrl, "https://")
		if isHTTP {
			transport := &mcp.StreamableClientTransport{
				Endpoint: url,
			}
			// Connect to a server over stdin/stdout.
			session, err = client.Connect(ctx, transport, nil)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, fmt.Errorf("error both cmd and url are emtpy")
	}

	defer session.Close()

	toolList, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &McpCommandHandler{
		name:  name,
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
		logrus.Debugf("trying to execute mcp remote server %s->%s", mcpServerName, mcpFunction)
		transport := &mcp.StreamableClientTransport{
			Endpoint: url,
		}
		// Connect to a server over stdin/stdout.
		session, err = client.Connect(ctx, transport, nil)
		if err != nil {
			return "", err
		}
	} else {
		logrus.Debugf("trying to execute mcp command %s->%s", mcpServerName, mcpFunction)
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

// GetMCPAuthPRM
// gets the protected resources metadata discovery
func GetMCPAuthPRM(url string) (string, error) {
	baseUrl, err := utils.GetBaseUrl(url)
	if err != nil {
		return "", fmt.Errorf("cannot get baseurl from url %s: %v", url, err)
	}
	prmUrl := baseUrl + "/.well-known/oauth-protected-resource/mcp"
	logrus.Infof("prmUrl %s", prmUrl)
	return utils.GetDataFromUrl(prmUrl)

}
