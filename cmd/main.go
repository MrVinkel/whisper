package main

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/mr_vinkel/whisper/cmd/whisper"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "whisper",
		Short: "Whisper secrets to your development environment",
	}

	rootCmd.AddCommand(whisper.SecretsCmd())
	rootCmd.AddCommand(whisper.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
