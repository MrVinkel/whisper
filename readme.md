# Whisper

Whisper secrets to your development environment

Whisper connects to a [HashiCorp Vault](https://www.vaultproject.io/) and uses [DirEnv](https://direnv.net/) to export secrets to our development environment as environment variables. This is useful for distributing secrets for test environments.

Whisper uses the developers own credentials and permission for fetching secrets. [DirEnv](https://direnv.net/) ensures the secrets are only loaded when entering a folder and is unloaded again when leaving the folder or closing the terminal. This avoids having random secrets floating around in files.

Whisper can also exec another program with the secrets set in the environment variables.

NOTE: Do not use whisper for using and/or distributing production secrets.

## Usage

```txt
Whisper secrets to your development environment

Usage:
  whisper [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  direnv      Whisper secrets to your development environment
  exec        Whisper secrets to an executable
  help        Help about any command
  version     Print version of whisper

Flags:
  -h, --help   help for whisper

Use "whisper [command] --help" for more information about a command.
```

See the [`test`](./test) folder for examples

### Configuration

See how to configure Vault for either [userpass](https://developer.hashicorp.com/vault/docs/auth/userpass) or [oidc](https://developer.hashicorp.com/vault/docs/auth/jwt)

Whisper requires a configuration file called `.whisper.yml`. This files configures how to authenticate to vault and which secrets to fetch. It should be placed in the root of a repository.

All whisper commands uses this file for fetching secrets

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
  # OIDC creates a callback to localhost:8250 by default. 
  # http://localhost:8250/oidc/callback should be configured as an allowed redirect uri in vault oidc and for the idp provider
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
  direnv_load whisper direnv
else
  echo "Please install whisper: https://github.com/mrvinkel/whisper"
fi
```

### Exec

Use `whisper exec` to execute a program with secrets defined in `.whisper.yml`.

Example:

```bash
# Load secrets into environment variables and execute another application
whisper exec -- node app.js
```

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
