package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// ClientConfig holds the overall configuration.
type ClientConfig struct {
	BaseURL  string `koanf:"baseURL"`
	ApiToken string `koanf:"apiToken"`
}

// validatePermissions verifies that provided file exists and has 0600 permissions
func validatePermissions(filePath string) error {
	// Clean the file path
	cleanPath := filepath.Clean(filePath)

	// Check if the file exists
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cannot read client configuration file %s: file does not exist", cleanPath)
		}
		return fmt.Errorf("cannot read client configuration file %s: %v", cleanPath, err)
	}

	// Skip permission check on Windows
	if runtime.GOOS == "windows" {
		return nil
	}

	// Retrieve the file mode
	fileMode := fileInfo.Mode()

	// Check if it's a regular file
	if !fileMode.IsRegular() {
		return fmt.Errorf("cannot read client configuration file %s: not a regular file", cleanPath)
	}

	// Extract the permission bits
	perm := fileMode.Perm()

	// Define the expected permission
	const expectedPerm os.FileMode = 0600

	// Compare the permissions
	if perm != expectedPerm {
		return fmt.Errorf("cannot read client configuration file %s: permissions are too open. It should be set to 0600", cleanPath)
	}

	return nil
}

// GetDefaultConfigPath returns path to the default client configuration file
func GetDefaultConfigPath() (string, error) {
	// Check if XDG_CONFIG_HOME is set
	var configPath string
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		configPath = filepath.Join(xdgConfigHome, "donezo", "config.yaml")
	} else {
		// Fallback to ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to determine user home directory: %v", err)
		}
		configPath = filepath.Join(homeDir, ".config", "donezo", "config.yaml")
	}
	return configPath, nil
}

// LoadClientConfig loads the configuration from the specified file path.
func LoadClientConfig(path string) (*ClientConfig, error) {
	var cfg ClientConfig
	if err := validatePermissions(path); err != nil {
		return nil, err
	}
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, err
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
