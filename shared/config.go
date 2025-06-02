package shared

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveConfig saves a configuration to a file
func SaveConfig(config interface{}, filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}

// LoadConfig loads a configuration from a file
func LoadConfig(filename string, config interface{}) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal JSON
	return json.Unmarshal(data, config)
}

// DefaultClientConfig creates a default client configuration
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		ServerAddress: "http://localhost:8080",
		ClientName:    "NetworkMonitor Client",
		Targets: []Target{
			{
				Name:     "Google",
				URL:      "https://www.google.com",
				Interval: 60,
				Enabled:  true,
			},
			{
				Name:     "GitHub",
				URL:      "https://github.com",
				Interval: 60,
				Enabled:  true,
			},
		},
		LogLevel: "info",
	}
}

// GetConfigDir returns the configuration directory for the application
func GetConfigDir(appName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".config", appName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetConfigFilePath returns the path to a config file
func GetConfigFilePath(appName, filename string) (string, error) {
	configDir, err := GetConfigDir(appName)
	if err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, filename), nil
}