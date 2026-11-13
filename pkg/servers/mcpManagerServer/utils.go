package mcpManagerServer

import (
	"fmt"

	pbMcp "github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/sirupsen/logrus"
)

const (
	SETTINGS_DIRECTORY = "configs"
	MCP_CONFIG_FILE    = "mcp-servers.config.yaml"
)

// Save
func save(data *pbMcp.MCPServerConfigs) error {

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, MCP_CONFIG_FILE)
	// Save list of configs
	if err := datalayer.SaveProtoYAML(filePath, data); err != nil {
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
func Load() (*pbMcp.MCPServerConfigs, error) {

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, MCP_CONFIG_FILE)
	data := &pbMcp.MCPServerConfigs{}

	// note: only the proto msg has to be sent, not the wrapper
	err := datalayer.ReadProtoYAML(filePath, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// findByUUIDName
// search MCP servers in config by UUID or name
// it returns the entry which matches either UUID or name
func findByUUIDName(uuid string, name string) (*pbMcp.MCPServerConfig, error) {

	entries, err := Load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
		return nil, err
	}

	for _, entry := range entries.Configs {
		if uuid == entry.Uuid {
			if uuid != "" && uuid == entry.Uuid {
				return entry, nil
			}
			if name != "" && name == entry.Name {
				return entry, nil
			}
		}
	}

	return nil, fmt.Errorf("mcp server with uuid/name %s/%s not found", uuid, name)

}

// getIndexByUUIDName
// search MCP servers in config by UUID or name
// it returns the index which matches either UUID or name
func getIndexByUUIDName(uuid string, name string) (*pbMcp.MCPServerConfigs, int, error) {

	entries, err := Load()
	if err != nil {
		logrus.Errorf("cannot read mcp server details - %v", err)
		return nil, -1, err
	}

	for index, entry := range entries.Configs {
		if uuid == entry.Uuid {
			if uuid != "" && uuid == entry.Uuid {
				return entries, index, nil
			}
			if name != "" && name == entry.Name {
				return entries, index, nil
			}
		}
	}

	return entries, -1, fmt.Errorf("mcp server with uuid/name %s/%s not found", uuid, name)

}
