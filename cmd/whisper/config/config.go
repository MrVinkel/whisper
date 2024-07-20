package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ENV_VAR_DIR_CONFIG = "WHISPER_DIR_CONFIG"
)

type KeyConfig struct {
	Name   string  `yaml:"name"`
	Rename *string `yaml:"rename,omitempty"`
}

type SecretConfig struct {
	Path   string      `yaml:"path"`
	Prefix *string     `yaml:"prefix,omitempty"`
	Keys   []KeyConfig `yaml:"keys"`
}

type DirConfig struct {
	Provider map[string]interface{} `yaml:"provider"`
	Secrets  []SecretConfig         `yaml:"secrets"`
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

	if config.Provider == nil {
		return nil, fmt.Errorf("provider is required")
	}
	if config.Provider["type"] == nil || config.Provider["type"] == "" {
		return nil, fmt.Errorf("provider type is required")
	}
	if config.Secrets == nil || len(config.Secrets) == 0 {
		return nil, fmt.Errorf("no secrets configured")
	}

	return config, nil
}

func (s *SecretConfig) Get(key string) *KeyConfig {
	for _, k := range s.Keys {
		if k.Name == key {
			return &k
		}
	}
	return nil
}

func (s *SecretConfig) Contains(key string) bool {
	for _, k := range s.Keys {
		if k.Name == key {
			return true
		}
	}
	return false
}
