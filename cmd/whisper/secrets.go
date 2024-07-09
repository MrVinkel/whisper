package whisper

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func SecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Whisper secrets to your development environment",
		Run:   Secrets,
	}

	cmd.Flags().BoolP("direnv", "", false, "export secrets to direnv by calling direnv dump")

	return cmd
}

func Secrets(cmd *cobra.Command, args []string) {
	config, err := ReadDirConfig()
	if err != nil {
		cmd.Printf("Failed to read config: %v\n", err)
		return
	}
	vault, err := Authenticate(config.Vault)
	if err != nil {
		cmd.Printf("Failed to authenticate: %v\n", err)
		return
	}
	secrets, err := vault.GetSecrets(config.Secrets)
	if err != nil {
		cmd.Printf("Failed to get secrets: %v\n", err)
		return
	}

	if ok, err := cmd.Flags().GetBool("direnv"); err == nil && ok {
		dump := exec.Command("direnv", "dump")
		dump.Env = os.Environ()
		for k, v := range secrets {
			dump.Env = append(dump.Env, fmt.Sprintf("%s=%s", k, v))
		}
		err := dump.Run()
		if err != nil {
			cmd.Printf("Failed to run direnv dump: %v\n", err)
		}
	}
}
