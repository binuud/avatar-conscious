package settings

import (
	"cmp"
	"maps"
	"slices"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/spf13/cobra"
)

func NewSettingsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "settings",
		Short: " *list*, *set*, *get*",
		Args:  cobra.MaximumNArgs(1),
	}

	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "lists all global settings",
		Run: func(cmd *cobra.Command, args []string) {
			if Global == nil {
				cmd.PrintErrln("cannot list settings, check if dataVolumes is mounted")
				return
			}
			mySettings := Global.GetAllByGroup()
			if mySettings == nil {
				cmd.PrintErrln("cannot list settings, check if dataVolumes is mounted")
				return
			}
			cmd.Println("Setting...")

			// sort map by group names
			// so display can be ordered
			sortedGroupNames := slices.Sorted(maps.Keys(mySettings))
			for _, groupName := range sortedGroupNames {
				cmd.Printf(" * **%s**\n", groupName)
				config := mySettings[groupName]

				// sort the config entries
				slices.SortFunc(config.Configs, func(a, b *pb.Config) int {
					return cmp.Compare(a.Name, b.Name)
				})
				for _, item := range config.Configs {
					cmd.Printf(" * * **%s**, %s\n", item.Name, item.Value)
				}
			}
		},
	})

	command.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "set a setting by name, set [name] [value]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				cmd.PrintErrln("please enter the name and value of the setting to be changed")
			}
			err := Global.Update(args[0], args[1])
			if err != nil {
				cmd.PrintErrln("cannot update settings entry")
			} else {
				cmd.Println("updated entry")
			}

		},
	})

	command.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "get a setting by name, get [name]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.PrintErrln("please enter the name of the setting to be displayed")
			}
			entry := Global.Get(args[0])
			if entry == nil {
				cmd.PrintErrln("cannot get settings entry")
			} else {
				cmd.Printf(" **%s** : %s\n", entry.Name, entry.Value)
				cmd.Printf(" %s\n", entry.Description)
			}

		},
	})

	return command
}
