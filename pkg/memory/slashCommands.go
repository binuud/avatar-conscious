package memory

import (
	"github.com/spf13/cobra"
)

func NewMemoryCommand() *cobra.Command {
	memoryCommand := &cobra.Command{
		Use:   "memory",
		Short: " *list*, *print*, *delete*",
		Args:  cobra.MaximumNArgs(1),
	}

	memoryCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "lists all memory files",
		Run: func(cmd *cobra.Command, args []string) {
			memories, err := ListMemory()
			if err != nil {
				cmd.PrintErrf("cannot list memories, check if dataVolumes is mounted, %v \n", err)
				return
			}
			cmd.Println("Memory files in the system...")
			for _, item := range memories {
				cmd.Println(" * ", item)
			}
		},
	})

	memoryCommand.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "deletes the specified memory",
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if name != "" {
				err := DeleteMemory(name)
				if err != nil {
					cmd.PrintErrf("cannot delete memory, '%s' %v \n", name, err)
				} else {
					cmd.Printf("Deleted memory %s\n", name)
				}
			}
		},
	})

	memoryCommand.AddCommand(&cobra.Command{
		Use:   "print",
		Short: "print the contents of the memory",
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if name == "" {
				cmd.PrintErrln("please enter the memory name, to view the contents")
				return
			}
			memory := NewMemory(name)
			memory.DumpMemory(cmd.OutOrStdout())
		},
	})

	return memoryCommand
}
