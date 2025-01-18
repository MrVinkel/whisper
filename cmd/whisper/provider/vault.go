package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/mrvinkel/whisper/cmd/whisper/config"
	"github.com/mrvinkel/whisper/cmd/whisper/util"
	"golang.org/x/term"
)

type VaultConfig struct {
	Address      string  `yaml:"address"`
	AuthMethod   string  `yaml:"authMethod"`
	AuthMount    *string `yaml:"authMount,omitempty"`
	SecretMount  string  `yaml:"secretMount"`
	CallbackPort *int    `yaml:"callbackPort,omitempty"`
}

type Vault struct {
	client *vault.Client
	config VaultConfig
}

func NewVaultProvider(config map[string]interface{}) (Provider, error) {
	var vaultConfig VaultConfig
	err := mapstructure.Decode(config, &vaultConfig)
	if err != nil {
		return nil, err
	}

	client, err := vault.New(
		vault.WithAddress(vaultConfig.Address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Vault{client: client, config: vaultConfig}, nil
}

func (v *Vault) Authenticate(ctx context.Context) error {
	switch v.config.AuthMethod {
	case "oidc":
		return v.oidc(ctx)
	case "userpass":
		return v.userpass(ctx)
	default:
		return fmt.Errorf("unsupported auth method: %s", v.config.AuthMethod)
	}
}

func (v *Vault) GetSecrets(secretConfig []config.SecretConfig) (map[string]string, error) {
	result := make(map[string]string)
	for _, secrets := range secretConfig {
		kvV2resp, err := v.client.Secrets.KvV2Read(context.Background(), secrets.Path, vault.WithMountPath(v.config.SecretMount))
		if err != nil {
			return nil, fmt.Errorf("failed to read secret: %w", err)
		}

		for k, v := range kvV2resp.Data.Data {
			keyConfig := secrets.Get(k)
			if len(secrets.Keys) != 0 && keyConfig == nil {
				continue
			}

			key := k
			value := fmt.Sprintf("%v", v)
			if keyConfig != nil && keyConfig.Rename != nil {
				key = *keyConfig.Rename
			} else if secrets.Prefix != nil {
				key = fmt.Sprintf("%s%s", *secrets.Prefix, k)
			}
			result[key] = value
		}
	}

	return result, nil
}

func (v *Vault) oidc(ctx context.Context) error {
	// Start callback server
	port := v.config.CallbackPort
	if port == nil {
		port = util.Ptr(8250)
	}
	callback := StartCallbackServer(*port)
	defer callback.Stop(ctx)

	// Make oidc auth request
	authMount := v.config.AuthMount
	if authMount == nil {
		authMount = util.Ptr("oidc")
	}
	r, err := v.client.Auth.JwtOidcRequestAuthorizationUrl(ctx,
		schema.JwtOidcRequestAuthorizationUrlRequest{
			RedirectUri: fmt.Sprintf("http://localhost:%d/oidc/callback", *port),
		},
		vault.WithMountPath(*authMount),
	)
	if err != nil {
		return err
	}

	// Open auth url in browser
	if u, ok := r.Data["auth_url"].(string); ok {
		if err := util.Open(u); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}
	} else {
		return fmt.Errorf("failed to get auth url")
	}

	// Wait for callback
	callbackURL := callback.WaitForCallback()

	// Handle callback
	r, err = v.client.Auth.JwtOidcCallback(ctx,
		"", // client nonce is not needed in this case
		callbackURL.Query().Get("code"),
		callbackURL.Query().Get("state"),
		vault.WithMountPath("oidc"),
	)
	if err != nil {
		return err
	}

	if err := v.client.SetToken(r.Auth.ClientToken); err != nil {
		return err
	}

	return nil
}

func (v *Vault) userpass(_ context.Context) error {
	fmt.Print("username: ")
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		return err
	}

	fmt.Print("password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println()

	mount := v.config.AuthMount
	if mount == nil {
		mount = util.Ptr("userpass")
	}

	resp, err := v.client.Auth.UserpassLogin(context.Background(), username, schema.UserpassLoginRequest{
		Password: string(password),
	}, vault.WithMountPath(*mount))
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	if err := v.client.SetToken(resp.Auth.ClientToken); err != nil {
		return err
	}

	return nil
}
