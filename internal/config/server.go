package config

import (
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// ServerConfig holds the overall configuration.
type ServerConfig struct {
	Port     int       `koanf:"port"`
	Database string    `koanf:"database"`
	JWT      JWTConfig `koanf:"jwt"`
}

type JWTConfig struct {
	Secret   []byte        `koanf:"secret"`
	Duration time.Duration `koanf:"duration"`
}

// LoadServerConfig loads the configuration from the specified file path.
func LoadServerConfig(path string) (ServerConfig, error) {
	var cfg ServerConfig
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return cfg, err
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
