package tools

import (
	"strings"

	"github.com/binuud/avatar-conscious/pkg/tools/aiGenTools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// This is a placeholder only
// actual implementation is in the command router of
// conscious server
// #TODO Binu this is a hack to speed up tool choosing
// else 2-3 seconds is spent on tool choosing by router
// on local inteferece
func NewImageAiCommand() *cobra.Command {

	command := &cobra.Command{
		Use:   "image",
		Short: "*create*, *edit*, *describe*",
		Args:  cobra.MaximumNArgs(1),
	}

	command.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "create an ai image with the given prompt",
		Run: func(cmd *cobra.Command, args []string) {
			imagePrompt := strings.Join(args[0:], " ")
			cmd.Printf("Creating image... %v", imagePrompt)
			imageFileName, err := aiGenTools.GetAiImage(imagePrompt, "")
			if err != nil {
				logrus.Errorf("err RouteOpenAiImageGenerationRequest -> GetAiImage result(%s), %v", imageFileName, err)
				cmd.PrintErrf("err RouteOpenAiImageGenerationRequest -> GetAiImage result(%s), %v\n", imageFileName, err)
			}
			cmd.Println(imageFileName)
		},
	})

	return command
}
