package whisper

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ENV_VAR_DIR_CONFIG = "WHISPER_DIR_CONFIG"
)

type SecretConfig struct {
	Path      string   `yaml:"path"`
	MountPath string   `yaml:"mount"`
	Prefix    string   `yaml:"prefix"`
	Keys      []string `yaml:"keys"`
}

type VaultConfig struct {
	Address      string `yaml:"address"`
	AuthMethod   string `yaml:"authMethod"`
	AuthMount    string `yaml:"authMount"`
	CallbackPort int    `yaml:"callbackPort"`
}

type DirConfig struct {
	Vault   VaultConfig    `yaml:"vault"`
	Secrets []SecretConfig `yaml:"secrets"`
}

func ReadDirConfig() (*DirConfig, error) {
	configPath := os.Getenv(ENV_VAR_DIR_CONFIG)
	if configPath == "" {
		configPath = ".whisper.yml"
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	content = []byte(os.ExpandEnv(string(content)))

	config := &DirConfig{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
