package avatarServer

import (
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

	// avatarCommand.AddCommand(&cobra.Command{
	// 	Use:   "delete",
	// 	Short: "deletes an avatar",
	// 	Run: func(cmd *cobra.Command, args []string) {

	// 		cmd.Printf("Deleting avatar, not yet implemented\n")

	// 	},
	// })

	// avatarCommand.AddCommand(&cobra.Command{
	// 	Use:   "start",
	// 	Short: "starts an avatar",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		cmd.Printf("Starting avatar - not implemented via command, use UI\n")
	// 	},
	// })

	avatarCommand.AddCommand(&cobra.Command{
		Use: "print",
		Short: `prints the steps and instructions of the AI Avatar
			print [avatar-name] : prints the steps in the avatar, build the avatar before printing`,
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
		Use: "build",
		Short: `build's the avatar. LLM is used to find the steps invovled in running the avatar.\n
			build [avatar-name] : builds the selected avatar (long waiting call)\n 
			build ALL : builds all avatars`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Println("wrong arguments, please refer to help ")
				return
			}
			name := args[0]
			if name == "ALL" {
				cmd.Println("building all items...")
				items, err := ListAvatars()
				if err != nil {
					cmd.PrintErrf("cannot list AI Avatars, check if dataVolumes is mounted, %v", err)
					return
				}

				for _, item := range items {
					avatarName := item[:len(item)-7]
					cmd.Println(" * Building Avatar ", avatarName)

					a := NewAvatarFromStore(avatarName)
					err := a.Build()
					if err != nil {
						cmd.PrintErrf("err failed %v \n", err)
					} else {
						cmd.Printf("Build completed for avatar %s\n", avatarName)
					}
				}
				return
			}
			if name != "" {
				cmd.Printf("Building avatar %s\n", name)
				a := NewAvatarFromStore(name)
				err := a.Build()
				if err != nil {
					cmd.PrintErrf("err failed %v \n", err)
				} else {
					cmd.Printf("Build completed for avatar %s\n", name)
				}
			}

		},
	})

	return avatarCommand
}
