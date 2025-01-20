package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func DirEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "direnv",
		Short: "Whisper secrets to your development environment",
		RunE:  DirEnv,
	}

	return cmd
}

func DirEnv(cmd *cobra.Command, args []string) error {
	secrets, err := readSecrets(cmd.Context())
	if err != nil {
		return err
	}

	dump := exec.Command("direnv", "dump")
	dump.Env = os.Environ()
	for k, v := range secrets {
		dump.Env = append(dump.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if err := dump.Run(); err != nil {
		return fmt.Errorf("failed to run direnv dump: %w", err)
	}
	return nil
}
