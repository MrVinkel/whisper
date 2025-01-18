# Whisper

Whisper secrets to your development environment

Whisper connects to a [HashiCorp Vault](https://www.vaultproject.io/) and uses [DirEnv](https://direnv.net/) to export secrets to our development environment as environment variables. This is useful for distributing secrets for test environments.

Whisper uses the developers own credentials and permission for fetching secrets. [DirEnv](https://direnv.net/) ensures the secrets are only loaded when entering a folder and is unloaded again when leaving the folder or closing the terminal. This avoids having random secrets floating around in files.

NOTE: Do not use whisper for using and/or distributing production secrets.

## Usage

See the [`test`](./test) folder for examples

### Configuration

See how to configure Vault for either [userpass](https://developer.hashicorp.com/vault/docs/auth/userpass) or [oidc](https://developer.hashicorp.com/vault/docs/auth/jwt)

Whisper uses per repository configuration called `.whisper.yml` to configure where to fetch secrets from and which secrets to fetch.

```yaml
provider:
  type: vault
  # Address to vault
  address: http://my-vault:8200

  # userpass authentication https://developer.hashicorp.com/vault/docs/auth/userpass
  # authMethod: userpass
  # authMount: userpass

  # OIDC authentication https://developer.hashicorp.com/vault/docs/auth/jwt
  authMethod: oidc
  authMount: oidc
  # OIDC creates a callback to localhost:8250 by default
  callbackPort: 8250

  # KV V2 mount to read secrets from
  secretMount: secret

# List of secrets to load
secrets:
  # Path to secret to load. All key values will be exported
- path: path/to/secret
  # Optional prefix for secrets
  prefix: MY_APP_

```

### DirEnv

Whisper uses [DirEnv](https://direnv.net/) for export secrets to environment variables

```bash
#!/bin/bash
if whisper version &>/dev/null; then
  direnv_load whisper secrets --direnv
else
  echo "Please install whisper: https://github.com/mrvinkel/whisper"
fi
```

### Exec

TODO

### DevBox

Whisper can be installed with [DevBox](https://www.jetify.com/devbox)

```json
{
  "$schema": "https://raw.githubusercontent.com/jetify-com/devbox/0.13.7/.schema/devbox.schema.json",
  "packages": [
    "github:mrvinkel/whisper"
  ]
}
```
