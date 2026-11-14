package mcpManagerServer

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

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

	instanceDetails, err := s.getInstanceDetails(in.Config)
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

	entries, err := Load()
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
		return nil, fmt.Errorf("mcp server entry already registered with avatar having uuid - %s", entries.Configs[found].Uuid)
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
	entries, err := Load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
		return nil, err
	}

	return &pbMcp.ListResponse{
		Servers: entries.Configs,
	}, nil
}

func (s *mcpManagerGrpcServer) Toggle(ctx context.Context, in *pbMcp.ToggleRequest) (*pbMcp.ToggleResponse, error) {

	logrus.Infof("Toggle called on MCP Manager %v", in.Uuid)

	entries, foundIndex, err := getIndexByUUIDName(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}

	if foundIndex != -1 {
		entries.Configs[foundIndex].Enabled = !entries.Configs[foundIndex].Enabled
		save(entries)
		return &pbMcp.ToggleResponse{
			Config: entries.Configs[foundIndex],
		}, nil
	}
	return nil, err
}

func (s *mcpManagerGrpcServer) Get(ctx context.Context, in *pbMcp.GetRequest) (*pbMcp.GetResponse, error) {

	logrus.Infof("Get called on MCP Manager %v", in)
	entry, err := findByUUIDName(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}

	instanceDetails, err := s.getInstanceDetails(entry)

	return &pbMcp.GetResponse{
		Instance: instanceDetails,
	}, nil

}

func (s *mcpManagerGrpcServer) Delete(ctx context.Context, in *pbMcp.DeleteRequest) (*pbMcp.DeleteResponse, error) {

	logrus.Infof("Delete called on MCP Manager %v", in)

	entries, foundIndex, err := getIndexByUUIDName(in.Uuid, "")
	if err != nil {
		logrus.Errorf("error - %v", err)
		return nil, err
	}

	if foundIndex == -1 {
		return nil, fmt.Errorf("mcp server with uuid %s, not found", in.Uuid)
	} else {
		entries.Configs = slices.Delete(entries.Configs, foundIndex, 1)
		save(entries)
	}

	return &pbMcp.DeleteResponse{
		Success: true,
	}, nil
}

func (s *mcpManagerGrpcServer) HandleTool(ctx context.Context, in *pbMcp.HandleToolRequest) (*pbMcp.HandleToolResponse, error) {

	logrus.Infof("HandleTool called on MCP Manager %v", in)
	entry, err := findByUUIDName(in.Uuid, "")
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

func (s *mcpManagerGrpcServer) getInstanceDetails(config *pbMcp.MCPServerConfig) (*pbMcp.MCPServerInstance, error) {

	mcp, err := mcpHandler.NewMcpCommandHandler(config.Name, config.Command, config.Url)
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

	// read the Protected Resource Metadata (PRM) document
	prmDoc := ""
	if config.McpType == pbMcp.MCPTransportType_MCP_TYPE_HTTP {
		prmDoc, err = mcpHandler.GetMCPAuthPRM(config.Url)
		if err != nil {
			prmDoc = fmt.Sprintf("error in fetching PRM Doc %v", err)
		} else {
			// try to pretty the json for display purpose
			prettyPrmDoc, err := utils.GetPrettyJson(prmDoc)
			if err == nil {
				// if no error, we send the prettyfied json back to client
				prmDoc = prettyPrmDoc
			}
		}
	}
	logrus.Infof("prmdoc %s", prmDoc)
	return &pbMcp.MCPServerInstance{
		Config: config,
		Tools:  mcpTools,
		PrmDoc: prmDoc,
	}, nil
}
