package mcpManagerServer

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"slices"

	pbMcp "github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/utils"
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
	mcp, err := mcpHandler.NewMcpCommandHandler(in.Config.Name, in.Config.Command, in.Config.Url)
	if err != nil {
		logrus.Errorf("cannot check mcp server, unreachable or command not found  - %v", err)
		return nil, fmt.Errorf("error -  checking mcp server, unreachable or command not found  - %v", err)
	}

	mcpTools := make([]*pbMcp.MCPTool, 0)
	for _, tool := range mcp.Tools {

		jsonBytes, err := json.MarshalIndent(tool.InputSchema, "", "  ")
		if err != nil {
			logrus.Errorf("error parse input schema %v", tool.InputSchema)
		}

		mcpTools = append(mcpTools, &pbMcp.MCPTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: string(jsonBytes),
		})
	}

	return &pbMcp.CheckResponse{
		Instance: &pbMcp.MCPServerInstance{
			Config: in.Config,
			Tools:  mcpTools,
		},
	}, nil
}

// Register
// check if endpoint is already register, if so flag an error
// if new mcp server, save it in list after checking reachability
func (s *mcpManagerGrpcServer) Register(ctx context.Context, in *pbMcp.RegisterRequest) (*pbMcp.RegisterResponse, error) {

	logrus.Infof("Register called on MCP Manager %v", in)

	entries, err := load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
	}

	if entries == nil {
		entries = &pbMcp.MCPServerConfigs{
			Configs: make([]*pbMcp.MCPServerConfig, 0),
		}
	}

	found := -1
	for i, entry := range entries.Configs {
		if in.Config.Command == entry.Command && in.Config.Url == entry.Url {
			found = i
			break
		}
	}

	if found != -1 {
		return nil, fmt.Errorf("mcp server entry already registered with avatar - %s", in.Config.Uuid)
	}

	in.Config.Uuid = utils.GenerateUUID()
	entries.Configs = append(entries.Configs, in.Config)

	slices.SortFunc(entries.Configs, func(a, b *pbMcp.MCPServerConfig) int {
		return cmp.Compare(a.Name, b.Name)
	})

	save(entries)

	return &pbMcp.RegisterResponse{
		Instance: &pbMcp.MCPServerInstance{
			Config: in.Config,
		},
	}, nil

}

// List
// Read the mcp config file and return all entries
func (s *mcpManagerGrpcServer) List(ctx context.Context, in *pbMcp.ListRequest) (*pbMcp.ListResponse, error) {

	logrus.Infof("List called on MCP Manager %v", in)
	entries, err := load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
		return nil, err
	}

	return &pbMcp.ListResponse{
		Servers: entries.Configs,
	}, nil
}

func (s *mcpManagerGrpcServer) Get(ctx context.Context, in *pbMcp.GetRequest) (*pbMcp.GetResponse, error) {

	logrus.Infof("Get called on MCP Manager %v", in)
	entries, err := load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
		return nil, err
	}

	found := -1
	for i, entry := range entries.Configs {
		if in.Uuid == entry.Uuid {
			found = i
			break
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("error mcp server with uuid %s not found", in.Uuid)
	}

	entry := entries.Configs[found]

	mcp, err := mcpHandler.NewMcpCommandHandler(entry.Name, entry.Command, entry.Url)
	if err != nil {
		logrus.Errorf("cannot check mcp server, unreachable or command not found  - %v", err)
		return nil, fmt.Errorf("error -  checking mcp server, unreachable or command not found  - %v", err)
	}

	mcpTools := make([]*pbMcp.MCPTool, 0)
	for _, tool := range mcp.Tools {

		jsonBytes, err := json.MarshalIndent(tool.InputSchema, "", "  ")
		if err != nil {
			logrus.Errorf("error parse input schema %v", tool.InputSchema)
		}

		mcpTools = append(mcpTools, &pbMcp.MCPTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: string(jsonBytes),
		})
	}

	return &pbMcp.GetResponse{
		Instance: &pbMcp.MCPServerInstance{
			Config: entry,
			Tools:  mcpTools,
		},
	}, nil

}

func (s *mcpManagerGrpcServer) Delete(ctx context.Context, in *pbMcp.DeleteRequest) (*pbMcp.DeleteResponse, error) {

	logrus.Infof("Delete called on MCP Manager %v", in)
	entries, err := load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
		return nil, err
	}

	found := -1
	for i, entry := range entries.Configs {
		if in.Uuid == entry.Uuid {
			found = i
			break
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("mcp server with uuid %s, not found", in.Uuid)
	} else {
		entries.Configs = slices.Delete(entries.Configs, found, 1)
		save(entries)
	}

	return &pbMcp.DeleteResponse{
		Success: true,
	}, nil
}
