package whisper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"golang.org/x/term"
)

type Vault struct {
	*vault.Client
}

func Authenticate(ctx context.Context, config VaultConfig) (*Vault, error) {
	if config.AuthMethod == "userpass" {
		return userpass(config)
	} else if config.AuthMethod == "azure" {
		return azure(config)
	} else if config.AuthMethod == "oidc" {
		return oidc(ctx, config)
	}
	return nil, fmt.Errorf("unsupported auth method: %s", config.AuthMethod)
}

type callback struct {
	done        chan bool
	callBackURL *url.URL
	server      http.Server
}

func startCallbackServer(port int) *callback {
	callback := &callback{
		done:        make(chan bool, 1),
		callBackURL: nil,
	}
	httpServer := &http.Server{Addr: fmt.Sprintf("localhost:%d", port)}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/oidc/callback", callback.handleCallback)
	httpServer.Handler = serverMux
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to start callback webserver: %s", err.Error())
		}
	}()
	return callback
}

func (c *callback) stop(ctx context.Context) {
	err := c.server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("failed to stop callback webserver: %v", err)
	}
}

func (c *callback) waitForCallback() *url.URL {
	<-c.done
	return c.callBackURL
}

func (c *callback) handleCallback(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, `
	<!DOCTYPE html>
	<html>
	<head>
	<script>
		setTimeout(function() {
			window.close()
		}, 2000);
	</script>
	</head>
	<body><p>Authenticated! This page closes in 2 seconds</p></body>
	</html>
	`)
	if err != nil {
		fmt.Printf("failed to write response: %v", err)
	}
	c.callBackURL = r.URL
	c.done <- true
}

func azure(config VaultConfig) (*Vault, error) {
	client, err := vault.New(
		vault.WithAddress(config.Address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Vault{client}, nil
}

func oidc(ctx context.Context, config VaultConfig) (*Vault, error) {
	client, err := vault.New(
		vault.WithAddress(config.Address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// Start callback server
	port := config.CallbackPort
	if port == 0 {
		port = 8250
	}
	callback := startCallbackServer(port)
	defer callback.stop(ctx)

	// Make oidc auth request
	authMount := config.AuthMount
	if authMount == "" {
		authMount = "oidc"
	}
	r, err := client.Auth.JwtOidcRequestAuthorizationUrl(ctx,
		schema.JwtOidcRequestAuthorizationUrlRequest{
			RedirectUri: fmt.Sprintf("http://localhost:%d/oidc/callback", port),
		},
		vault.WithMountPath(authMount),
	)
	if err != nil {
		return nil, err
	}

	// Open auth url in browser
	if u, ok := r.Data["auth_url"].(string); ok {
		if err := Open(u); err != nil {
			return nil, fmt.Errorf("failed to open browser: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to get auth url")
	}

	// Wait for callback
	callbackURL := callback.waitForCallback()

	// Handle callback
	r, err = client.Auth.JwtOidcCallback(ctx,
		"", // client nonce is not needed in this case
		callbackURL.Query().Get("code"),
		callbackURL.Query().Get("state"),
		vault.WithMountPath("oidc"),
	)
	if err != nil {
		return nil, err
	}

	if err := client.SetToken(r.Auth.ClientToken); err != nil {
		return nil, err
	}

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
