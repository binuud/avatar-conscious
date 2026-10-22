package testbed

import (
	"github.com/spf13/cobra"
)

func NewTestBedCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "testbed",
		Short: " *ui*",
		Args:  cobra.MaximumNArgs(1),
	}

	command.AddCommand(&cobra.Command{
		Use:   "ui",
		Short: "returns ui to render",
		Run: func(cmd *cobra.Command, args []string) {
			p := NewPane()
			text, err := p.ShowPane()
			if err != nil {
				cmd.PrintErrf("cannot list memories, check if dataVolumes is mounted, %v \n", err)
				return
			}
			cmd.Println(text)
		},
	})

	return command
}
