package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func ExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "exec",
		Short:                 "Whisper secrets to an executable",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Example:               "whisper exec -- node app.js",
		RunE:                  Exec,
	}
	return cmd
}

func Exec(cmd *cobra.Command, args []string) error {
	secrets, err := readSecrets(cmd.Context())
	if err != nil {
		return err
	}

	env := make([]string, 0, len(secrets))
	for k, v := range secrets {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	command := exec.Command(args[0], args[1:]...)
	command.Env = append(env, os.Environ()...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if errors.Is(command.Err, exec.ErrDot) {
		command.Err = nil
	}
	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to exec: %d - %w", command.ProcessState.ExitCode(), err)
	}
	return nil
}
