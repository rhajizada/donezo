package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config holds the overall configuration.
type Config struct {
	Port     int    `koanf:"port"`
	Database string `koanf:"database"`
}

// Load loads the configuration from the specified file path.
func Load(path string) (Config, error) {
	var cfg Config
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return cfg, err
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
