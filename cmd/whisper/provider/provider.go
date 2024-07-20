package provider

import (
	"context"
	"fmt"

	"gitlab.com/mr_vinkel/whisper/cmd/whisper/config"
)

type Provider interface {
	Authenticate(ctx context.Context) error
	GetSecrets(secrets []config.SecretConfig) (map[string]string, error)
}

func NewProvider(ctx context.Context, config map[string]interface{}) (Provider, error) {
	if config["type"] == "vault" {
		return NewVaultProvider(config)
	}
	return nil, fmt.Errorf("unsupported provider type: %s", config["type"])
}
