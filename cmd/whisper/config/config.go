package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ENV_VAR_CONFIG = "WHISPER_CONFIG"
)

var defaultConfigs = []string{
	".whisper.yml",
	".whisper.yaml",
	"whisper.yml",
	"whisper.yaml",
}

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
	configPath := os.Getenv(ENV_VAR_CONFIG)
	configs := defaultConfigs
	if configPath != "" {
		configs = append([]string{configPath}, configs...)
	}

	var err error
	var content []byte
	found := false
	for _, path := range configs {
		content, err = os.ReadFile(path)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no config file found")
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
	if len(config.Secrets) == 0 {
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
