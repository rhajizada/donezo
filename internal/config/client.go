package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// ClientConfig holds the overall configuration.
type ClientConfig struct {
	BaseURL  string `koanf:"baseURL"`
	ApiToken string `koanf:"apiToken"`
}

// LoadClientConfig loads the configuration from the specified file path.
func LoadClientConfig(path string) (ClientConfig, error) {
	var cfg ClientConfig
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return cfg, err
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
