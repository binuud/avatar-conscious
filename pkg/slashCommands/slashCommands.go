package slashCommands

// slashCommands
// Binu Udayakumar
//
// from the ai prompt, users can interact and debug the avatar
// each package has to implement a slashCommands.go
// where the package implements the commands that users can debug
// or interact with the package

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/memory"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarServer"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/testbed"
	"github.com/spf13/cobra"
)

var SlashCommands = &cobra.Command{
	Use:     "/[command] [arg] [flags]",
	Short:   "All commands start with a slash, eg: /help",
	Long:    `Slash commands can be invoked from the chat input.`,
	Example: `  /avatar list`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true, // This completely removes the 'completion' subcommand
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
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

// Bootstrap()
// add all slash commands from the different packages
func Bootstrap() {
	SlashCommands.AddCommand(llm.NewModelsCommand())
	SlashCommands.AddCommand(avatarServer.NewAvatarCommand())
	SlashCommands.AddCommand(memory.NewMemoryCommand())
	SlashCommands.AddCommand(settings.NewSettingsCommand())
	SlashCommands.AddCommand(testbed.NewTestBedCommand())

	Execute()
}
