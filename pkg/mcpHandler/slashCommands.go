package mcpHandler

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewMCPServerCommand() *cobra.Command {

	avatarCommand := &cobra.Command{
		Use:   "mcp",
		Short: "*tool-list*, *call*",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Println("Available subcommands under mcp:")
			// Loop through all attached subcommands
			for _, subCmd := range cmd.Commands() {
				cmd.Printf("  - %s: %s\n", subCmd.Name(), subCmd.Short)
			}
		},
	}

	avatarCommand.AddCommand(&cobra.Command{
		Use: "tool-list [mcp-server-cli-command]",
		Short: `lists all the tools and description provided by the mcp server, which is runnable from CLI.\n
			The cli command should be in the present in the PATH\n `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.PrintErrln("wrong arguments, please refer to help ")
				return
			}
			mcpName := args[0]
			mcp, err := NewMcpCommandHandler(mcpName)
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

			outMsg, err := HandleMcp(mcpName, toolName, paramMap)
			if err != nil {
				cmd.PrintErrf("error calling tool %s / %s -  %v\n", mcpName, toolName, err)
				logrus.Errorf("tool call %v", err)
				return
			}
			cmd.Printf("Tool Output: \n%v\n", outMsg)

		},
	})

	return avatarCommand
}
