package whisper

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.com/mr_vinkel/whisper/cmd/whisper/config"
	"gitlab.com/mr_vinkel/whisper/cmd/whisper/provider"
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
	config, err := config.ReadDirConfig()
	if err != nil {
		cmd.Printf("Failed to read config: %v\n", err)
		return
	}
	secretProvider, err := provider.NewProvider(cmd.Context(), config.Provider)
	if err != nil {
		cmd.Printf("Failed to create provider: %v\n", err)
		return
	}
	err = secretProvider.Authenticate(cmd.Context())
	if err != nil {
		cmd.Printf("Failed to authenticate: %v\n", err)
		return
	}
	secrets, err := secretProvider.GetSecrets(config.Secrets)
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
