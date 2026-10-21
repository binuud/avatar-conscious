package llm

import "github.com/spf13/cobra"

func NewModelsCommand() *cobra.Command {
	modelsCommand := &cobra.Command{
		Use:   "models",
		Short: "list models",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "world"
			if len(args) == 1 {
				name = args[0]
			}
			cmd.Printf("Hello, %s!\n", name)
			return nil
		},
	}

	modelsCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "lists all models",
		Run: func(cmd *cobra.Command, args []string) {
			models, err := MakeLLMGetModels()
			if err != nil {
				cmd.Println("Cannot connect to llm server")
				return
			}
			cmd.Println("Model listing from llm server")
			for _, model := range models {
				cmd.Printf(" * %s\n", model.Id)
			}
		},
	})

	return modelsCommand
}
