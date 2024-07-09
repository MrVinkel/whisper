package whisper

import "github.com/spf13/cobra"

var (
	Version = "dev"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version of whisper",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("whisper %s\n", Version)
		},
	}
}
