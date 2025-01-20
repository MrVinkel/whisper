package cmd

import (
	"context"
	"fmt"

	"github.com/mrvinkel/whisper/cmd/whisper/config"
	"github.com/mrvinkel/whisper/cmd/whisper/provider"
)

func readSecrets(ctx context.Context) (map[string]string, error) {
	config, err := config.ReadDirConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	secretProvider, err := provider.NewProvider(ctx, config.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	err = secretProvider.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	secrets, err := secretProvider.GetSecrets(config.Secrets)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}
