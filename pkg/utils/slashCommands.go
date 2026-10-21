package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var SlashCommands = &cobra.Command{
	Use:     "/[command] [arg]",
	Short:   "All commands start with a slash, eg: /help",
	Long:    `Slash commands can be invoked from the chat input.`,
	Example: `  /avatar list --all`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			// os.Exit(0)
		}
		// do actual work
	},
}

func NewModelsCommand() *cobra.Command {
	modelsCommand := &cobra.Command{
		Use:   "greet [name]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "world"
			if len(args) == 1 {
				name = args[0]
			}
			fmt.Printf("Hello, %s!\n", name)
			return nil
		},
	}

	modelsCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "lists all models",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Models List\n")
		},
	})

	return modelsCommand
}

func NewAvatarCommand() *cobra.Command {
	avatarCommand := &cobra.Command{
		Use:   "avatar",
		Short: "*list*, *start*, *update*",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Println("Available subcommands under avatar:")
			// Loop through all attached subcommands
			for _, subCmd := range cmd.Commands() {
				cmd.Printf("  - %s: %s\n", subCmd.Name(), subCmd.Short)
			}
		},
	}

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "lists an avatar",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Listing avatar\n")
		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "deletes an avatar",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Printf("Deleting avatar\n")

		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "starts an avatar",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Printf("Starting avatar\n")

		},
	})

	return avatarCommand
}

func Execute() {
	if err := SlashCommands.Execute(); err != nil {
		os.Exit(1)
	}
}

func SetStream(buf *bytes.Buffer) {
	SlashCommands.SetOut(buf)
	SlashCommands.SetErr(buf)
}

func HandleSlashCommand(input string) {
	replArgs := strings.Fields(input)

	// 5. Intercept Cobra behavior for REPL context
	// Silence errors and usage printing so a typo doesn't kill or flood your loop
	SlashCommands.SilenceErrors = true
	SlashCommands.SilenceUsage = true

	// 6. Feed the arguments into Cobra
	SlashCommands.SetArgs(replArgs)

	// 7. Execute! If it returns an error, catch it and print it locally
	if err := SlashCommands.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func BootstrapSlashCommands() {
	SlashCommands.AddCommand(NewModelsCommand())
	SlashCommands.AddCommand(NewAvatarCommand())
	Execute()
}
