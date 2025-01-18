package main

import (
	"os"

	"github.com/mrvinkel/whisper/cmd/whisper/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "whisper",
		Short: "Whisper secrets to your development environment",
	}

	rootCmd.AddCommand(cmd.SecretsCmd())
	rootCmd.AddCommand(cmd.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
