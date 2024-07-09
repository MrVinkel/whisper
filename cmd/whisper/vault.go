package whisper

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"golang.org/x/term"
)

type Vault struct {
	*vault.Client
}

func Authenticate(config VaultConfig) (*Vault, error) {
	if config.AuthMethod == "userpass" {
		return userpass(config)
	} else if config.AuthMethod == "azure" {
		return azure(config)
	} else if config.AuthMethod == "oidc" {
		return oidc(config)
	}
	return nil, fmt.Errorf("unsupported auth method: %s", config.AuthMethod)
}

func azure(config VaultConfig) (*Vault, error) {
	client, err := vault.New(
		vault.WithAddress(config.Address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// https://blue42.net/code/oidc-login-cli-vault/post/

	return &Vault{client}, nil
}

func oidc(config VaultConfig) (*Vault, error) {
	client, err := vault.New(
		vault.WithAddress(config.Address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// https://blue42.net/code/oidc-login-cli-vault/post/

	return &Vault{client}, nil
}

func userpass(config VaultConfig) (*Vault, error) {
	client, err := vault.New(
		vault.WithAddress(config.Address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	fmt.Print("username: ")
	var username string
	_, err = fmt.Scanln(&username)
	if err != nil {
		return nil, err
	}

	fmt.Print("password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	fmt.Println()

	mount := config.AuthMount
	if mount == "" {
		mount = "userpass"
	}

	resp, err := client.Auth.UserpassLogin(context.Background(), username, schema.UserpassLoginRequest{
		Password: string(password),
	}, vault.WithMountPath(mount))
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		return nil, err
	}

	return &Vault{client}, nil
}

func (v *Vault) GetSecrets(configs []SecretConfig) (map[string]string, error) {
	result := make(map[string]string)
	for _, config := range configs {
		kvV2resp, err := v.Secrets.KvV2Read(context.Background(), config.Path, vault.WithMountPath(config.MountPath))
		if err != nil {
			return nil, fmt.Errorf("failed to read secret: %w", err)
		}

		for k, v := range kvV2resp.Data.Data {
			if len(config.Keys) > 0 && !slices.Contains(config.Keys, k) {
				continue
			}

			key := k
			value := fmt.Sprintf("%v", v)
			if config.Prefix != "" {
				key = fmt.Sprintf("%s%s", config.Prefix, k)
			}
			result[key] = value
		}
	}

	return result, nil
}
