package mcpManagerServer

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	pbMcp "github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/sirupsen/logrus"
)

const (
	SETTINGS_DIRECTORY = "configs"
	MCP_CONFIG_FILE    = "mcp-servers.config.yaml"
)

type mcpServers struct {
	Servers *pbMcp.MCPServerConfigs
}

func NewMcpServers() *mcpServers {
	mcpServers := &mcpServers{}
	err := mcpServers.Load()
	if err != nil {
		// probably list emtpy
		// handle scenario later
	}
	return mcpServers
}

// Save
func (s *mcpServers) save() error {

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, MCP_CONFIG_FILE)
	// Save list of configs
	if err := datalayer.SaveProtoYAML(filePath, s.Servers); err != nil {
		logrus.Errorf("Cannot save config file %s %v", filePath, err)
		return err
	}
	logrus.Infof("Saved MCP Configs to file %s", MCP_CONFIG_FILE)
	return nil
}

// Load
// Read the config file from storage
// Update the values from storage to the current settings
// Do not add entries, only updation to existing entries happens here
func (s *mcpServers) Load() error {

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, MCP_CONFIG_FILE)
	s.Servers = &pbMcp.MCPServerConfigs{}

	// note: only the proto msg has to be sent, not the wrapper
	err := datalayer.ReadProtoYAML(filePath, s.Servers)
	if err != nil {
		return err
	}

	return nil
}

// findByUUIDName
// search MCP servers in config by UUID or name
// it returns the entry which matches either UUID or name
func (s *mcpServers) FindByUUIDName(uuid string, name string) (*pbMcp.MCPServerConfig, error) {

	for _, entry := range s.Servers.Configs {

		if uuid != "" && uuid == entry.Uuid {
			return entry, nil
		}
		if name != "" && name == entry.Name {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("mcp server with uuid/name %s/%s not found", uuid, name)

}

// getIndexByUUIDName
// search MCP servers in config by UUID or name
// it returns the index which matches either UUID or name
func (s *mcpServers) GetIndexByUUIDName(uuid string, name string) (int, error) {

	for index, entry := range s.Servers.Configs {

		if uuid != "" && uuid == entry.Uuid {
			return index, nil
		}
		if name != "" && name == entry.Name {
			return index, nil
		}
	}

	return -1, fmt.Errorf("mcp server with uuid/name %s/%s not found", uuid, name)

}

// Register
// Register the new mcp server
// search by name, to see if the mcp server is already registered
func (s *mcpServers) Register(inConfig *pbMcp.MCPServerConfig) (*pbMcp.MCPServerConfig, error) {

	if s.Servers == nil {
		s.Servers = &pbMcp.MCPServerConfigs{
			Configs: make([]*pbMcp.MCPServerConfig, 0),
		}
	}

	// search the server list by command or url
	found := -1
	for i, entry := range s.Servers.Configs {
		if inConfig.Command == entry.Command && inConfig.Url == entry.Url {
			found = i
			break
		}
	}
	if found != -1 {
		return nil, fmt.Errorf("mcp server entry already registered with avatar having uuid - %s", s.Servers.Configs[found].Uuid)
	}

	inConfig.Uuid = utils.GenerateUUID()
	s.Servers.Configs = append(s.Servers.Configs, inConfig)

	slices.SortFunc(s.Servers.Configs, func(a, b *pbMcp.MCPServerConfig) int {
		return cmp.Compare(a.Name, b.Name)
	})

	s.save()

	return inConfig, nil
}

// Toggle
// search MCP servers in config by UUID or name
// then toggle and save the list
func (s *mcpServers) Toggle(uuid string, name string) (*pbMcp.MCPServerConfig, error) {

	found, err := s.GetIndexByUUIDName(uuid, name)
	if err != nil {
		return nil, err
	}

	s.Servers.Configs[found].Enabled = !s.Servers.Configs[found].Enabled
	s.save()
	return s.Servers.Configs[found], nil

}

// Delete
// search MCP servers in config by UUID or name
// then delete and save the list
func (s *mcpServers) Delete(uuid string, name string) error {

	found, err := s.GetIndexByUUIDName(uuid, name)
	if err != nil {
		return err
	}

	s.Servers.Configs = slices.Delete(s.Servers.Configs, found, 1)
	s.save()
	return nil

}

func (s *mcpServers) GetInstanceDetails(config *pbMcp.MCPServerConfig) (*pbMcp.MCPServerInstance, error) {

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
