package mcpManagerServer

import (
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
func load() (*pbMcp.MCPServerConfigs, error) {

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, MCP_CONFIG_FILE)
	data := &pbMcp.MCPServerConfigs{}

	// note: only the proto msg has to be sent, not the wrapper
	err := datalayer.ReadProtoYAML(filePath, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
