package mcpManagerServer

import (
	"encoding/json"
	"strings"

	"github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/mcpHandler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewMCPServerCommand() *cobra.Command {

	avatarCommand := &cobra.Command{
		Use:   "mcp",
		Short: "*servers*, *tools*, *call*",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Println("Available subcommands under mcp:")
			// Loop through all attached subcommands
			for _, subCmd := range cmd.Commands() {
				cmd.Printf("  - %s: %s\n", subCmd.Name(), subCmd.Short)
			}
		},
	}

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "servers",
		Short: `lists all enabled mcp servers.\n`,
		Run: func(cmd *cobra.Command, args []string) {

			mcpServers := NewMcpServers()

			for _, mcpServer := range mcpServers.Servers.Configs {
				cmd.Printf("* Name: %s, Enabled: %v\n", mcpServer.Name, mcpServer.Enabled)
				switch mcpServer.McpType {
				case mcpManager.MCPTransportType_MCP_TYPE_STDIO:
					cmd.Printf("* Command: %s\n", mcpServer.Command)
				case mcpManager.MCPTransportType_MCP_TYPE_HTTP:
					cmd.Printf("* Remote: %s\n", mcpServer.Url)
				}
				cmd.Println("--------------------------------------")
			}

		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use: "tools [mcp-server-cli-command]",
		Short: `lists all the tools and description provided by the mcp server, which is runnable from CLI.\n
			The cli command should be in the present in the PATH\n `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.PrintErrln("wrong arguments, please refer to help ")
				return
			}
			mcpName := args[0]

			mcpServers := NewMcpServers()
			entry, err := mcpServers.FindByUUIDName("", mcpName)
			if err != nil {
				cmd.PrintErrf("Cannot get tools from mcp server/commnad %s : %v\n", mcpName, err)
				return
			}

			mcp, err := mcpHandler.NewMcpCommandHandler(entry.Name, entry.Command, "")
			if err != nil {
				cmd.PrintErrf("Cannot get tools from mcp server/commnad %s : %v\n", mcpName, err)
				logrus.Errorf("Cannot get tools from mcp server/commnad %s : %v", mcpName, err)
			}

			// for each tool in mcp create an entry in GlobalTools
			for _, tool := range mcp.Tools {
				cmd.Printf("* %s\n", tool.Name)
				cmd.Printf("* * %s\n", tool.Description)

				// Marshal returns a []byte slice
				jsonBytes, err := json.MarshalIndent(tool.InputSchema, "", "  ")
				if err != nil {
					cmd.PrintErrf("error parse input schema %v\n", tool.InputSchema)
					logrus.Errorf("error parse input schema %v", tool.InputSchema)
				}
				cmd.Printf("* * Input Schema \n%v\n", string(jsonBytes))
			}

		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use: "call [mcp-server-cli-command]",
		Short: `calls a particular tool of the mcp server, which is runnable from CLI.\n
			The cli command should be in the present in the PATH\n `,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.PrintErrln("wrong arguments, please refer to help ")
				return
			}
			params := ""
			mcpName := args[0]
			toolName := args[1]

			if len(args) > 2 {
				params = strings.Join(args[2:], "")
			}
			cmd.Printf("tool name %s/%s: params %s\n", mcpName, toolName, params)

			var paramMap map[string]any

			// Unmarshal the JSON bytes into the map pointer
			err := json.Unmarshal([]byte(params), &paramMap)
			if err != nil {
				cmd.PrintErrf("error unmarshalling params %s : %v \n", params, err)
				return
			}

			mcp := NewMcpServers()
			outMsg, err := mcp.HandleTool("", mcpName, toolName, params)
			if err != nil {
				cmd.PrintErrf("error calling tool %s / %s -  %v\n", mcpName, toolName, err)
				return
			}

			cmd.Printf("Tool Output: \n%v\n", outMsg)

		},
	})

	return avatarCommand
}
