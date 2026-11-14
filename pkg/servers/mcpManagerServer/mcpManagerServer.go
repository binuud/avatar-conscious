package mcpManagerServer

import (
	"context"
	"fmt"
	"strings"

	pbMcp "github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/sirupsen/logrus"
)

type mcpManagerGrpcServer struct {
	pbMcp.UnimplementedMCPManagerServiceServer
}

func NewMCPManagerGrpcServer() *mcpManagerGrpcServer {
	return &mcpManagerGrpcServer{}
}

// Check
// check tools, reachability of the mcp server
func (s *mcpManagerGrpcServer) Check(ctx context.Context, in *pbMcp.CheckRequest) (*pbMcp.CheckResponse, error) {

	logrus.Infof("Check called on MCP Manager %v", in)
	mcp := NewMcpServers()
	instanceDetails, err := mcp.GetInstanceDetails(in.Config)
	if err != nil {
		return nil, err
	}

	return &pbMcp.CheckResponse{
		Instance: instanceDetails,
	}, nil

}

// Register
// check if endpoint is already register, if so flag an error
// if new mcp server, save it in list after checking reachability
func (s *mcpManagerGrpcServer) Register(ctx context.Context, in *pbMcp.RegisterRequest) (*pbMcp.RegisterResponse, error) {

	logrus.Infof("Register called on MCP Manager %v", in)

	if in.Config.Name == "" {
		return nil, fmt.Errorf("emtpy MCP Server name not accepted")
	}
	in.Config.Name = strings.ToLower(in.Config.Name)

	if in.Config.McpType == pbMcp.MCPTransportType_MCP_TYPE_HTTP && in.Config.Url == "" {
		return nil, fmt.Errorf("emtpy MCP Server Remote URL not accepted")
	}

	if in.Config.McpType == pbMcp.MCPTransportType_MCP_TYPE_STDIO && in.Config.Command == "" {
		return nil, fmt.Errorf("emtpy MCP Server Command not accepted")
	}

	mcp := NewMcpServers()
	config, err := mcp.Register(in.Config)
	if err != nil {
		return nil, err
	}

	// return the config with generated uuid
	return &pbMcp.RegisterResponse{
		Instance: &pbMcp.MCPServerInstance{
			Config: config,
		},
	}, nil

}

// List
// Read the mcp config file and return all entries
func (s *mcpManagerGrpcServer) List(ctx context.Context, in *pbMcp.ListRequest) (*pbMcp.ListResponse, error) {

	logrus.Infof("List called on MCP Manager %v", in)
	mcp := NewMcpServers()

	return &pbMcp.ListResponse{
		Servers: mcp.Servers.Configs,
	}, nil
}

func (s *mcpManagerGrpcServer) Toggle(ctx context.Context, in *pbMcp.ToggleRequest) (*pbMcp.ToggleResponse, error) {

	logrus.Infof("Toggle called on MCP Manager %v", in.Uuid)

	mcp := NewMcpServers()
	mcpConfig, err := mcp.Toggle(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}

	return &pbMcp.ToggleResponse{
		Config: mcpConfig,
	}, nil
}

func (s *mcpManagerGrpcServer) Get(ctx context.Context, in *pbMcp.GetRequest) (*pbMcp.GetResponse, error) {

	logrus.Infof("Get called on MCP Manager %v", in)
	mcp := NewMcpServers()
	entry, err := mcp.FindByUUIDName(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}

	instanceDetails, err := mcp.GetInstanceDetails(entry)

	return &pbMcp.GetResponse{
		Instance: instanceDetails,
	}, nil

}

func (s *mcpManagerGrpcServer) Delete(ctx context.Context, in *pbMcp.DeleteRequest) (*pbMcp.DeleteResponse, error) {

	logrus.Infof("Delete called on MCP Manager %v", in)

	mcp := NewMcpServers()
	err := mcp.Delete(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}

	return &pbMcp.DeleteResponse{
		Success: true,
	}, nil
}

func (s *mcpManagerGrpcServer) HandleTool(ctx context.Context, in *pbMcp.HandleToolRequest) (*pbMcp.HandleToolResponse, error) {

	logrus.Infof("HandleTool called on MCP Manager %v", in)
	mcp := NewMcpServers()
	entry, err := mcp.FindByUUIDName(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}
	logrus.Debugf("HandleTool MCP Server with name %s found", entry.Name)
	command := ""
	switch entry.McpType {
	case pbMcp.MCPTransportType_MCP_TYPE_STDIO:
		command = entry.Command
	case pbMcp.MCPTransportType_MCP_TYPE_HTTP, pbMcp.MCPTransportType_MCP_TYPE_SSE:
		command = entry.Url
	}
	outMsg, err := mcpHandler.HandleMcpTool(command, in.ToolName, in.InputData)
	if err != nil {
		return nil, err
	}

	return &pbMcp.HandleToolResponse{
		Output: outMsg,
	}, nil
}
