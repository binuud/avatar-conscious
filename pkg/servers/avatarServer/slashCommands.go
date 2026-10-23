package avatarServer

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAvatarCommand() *cobra.Command {
	avatarCommand := &cobra.Command{
		Use:   "avatar",
		Short: "*list*, *start*, *build*, *print*",
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
		Short: "lists all avatars",
		Run: func(cmd *cobra.Command, args []string) {
			items, err := ListAvatars()
			if err != nil {
				cmd.PrintErrf("cannot list AI Avatars, check if dataVolumes is mounted, %v", err)
				return
			}
			cmd.Println("AI Avatar...")
			for _, item := range items {
				cmd.Println(" * ", item[:len(item)-7])
			}
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

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "print",
		Short: "prints the steps and instructions of the AI Avatar",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				name := args[0]
				if name != "" {
					a := NewAvatar(name)
					a.Print(cmd.OutOrStdout())
				}
			}
		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "build",
		Short: "builds an avatar, building an avatar, saves times during running an avatar agent",
		Run: func(cmd *cobra.Command, args []string) {
			buildAll, _ := cmd.Flags().GetBool("all")
			if buildAll {
				fmt.Println("building all items...")
			} else if len(args) == 1 {
				name := args[0]
				if name != "" {
					cmd.Printf("Building avatar %s\n", name)
					a := NewAvatarFromStore("election-officer.md")
					a.Build()
					a.Save()
				}
			}
		},
	})

	return avatarCommand
}
