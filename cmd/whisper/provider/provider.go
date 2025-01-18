package provider

import (
	"context"
	"fmt"

	"github.com/mrvinkel/whisper/cmd/whisper/config"
)

type Provider interface {
	Authenticate(ctx context.Context) error
	GetSecrets(secrets []config.SecretConfig) (map[string]string, error)
}

func NewProvider(ctx context.Context, config map[string]interface{}) (Provider, error) {
	provider := config["type"]
	switch provider {
	case "vault":
		return NewVaultProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", provider)
	}
}
