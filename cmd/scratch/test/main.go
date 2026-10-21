package main

import (
	"fmt"
	"os"

	bootstap "github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "demo [arg]",
	Short:   "A brief description shown in command lists",
	Long:    `A detailed paragraph explaining what the command does when help is explicitly invoked.`,
	Example: `  demo my-argument --flag-value`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		// do actual work
	},
}

func NewGreetCommand() *cobra.Command {
	return &cobra.Command{
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
}

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().Int("port", 8080, "port to listen on")
	return cmd
}

func NewDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "debug the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().Int("port", 8080, "port to listen on")
	return cmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	// rootCmd.AddCommand(NewServeCommand())
	// rootCmd.AddCommand(NewDebugCommand())
	// rootCmd.AddCommand(NewGreetCommand())
	fmt.Println("Started")
	bootstap.Bootstrap()
	// Execute()
}
