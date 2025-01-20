package main

import (
	"fmt"
	"os"

	"github.com/mrvinkel/whisper/cmd/whisper/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "whisper",
		Short:         "Whisper secrets to your development environment",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.AddCommand(cmd.DirEnvCmd())
	rootCmd.AddCommand(cmd.VersionCmd())
	rootCmd.AddCommand(cmd.ExecCmd())

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErr(err)
		fmt.Println()
		fmt.Println()
		rootCmd.Usage() // nolint: errcheck
		os.Exit(1)
	}
}
