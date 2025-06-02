package client

import (
	"networkmonitor/shared"
	"path/filepath"
)

const (
	// AppName is the name of the application
	AppName = "NetworkMonitor"
	
	// ConfigFileName is the name of the config file
	ConfigFileName = "client.json"
)

// LoadClientConfig loads the client configuration
func LoadClientConfig() (shared.ClientConfig, error) {
	// Get config file path
	configPath, err := shared.GetConfigFilePath(AppName, ConfigFileName)
	if err != nil {
		return shared.DefaultClientConfig(), err
	}

	// Load config
	var config shared.ClientConfig
	if err := shared.LoadConfig(configPath, &config); err != nil {
		// If file doesn't exist, create default config
		config = shared.DefaultClientConfig()
		if err := SaveClientConfig(config); err != nil {
			return config, err
		}
	}

	return config, nil
}

// SaveClientConfig saves the client configuration
func SaveClientConfig(config shared.ClientConfig) error {
	// Get config file path
	configPath, err := shared.GetConfigFilePath(AppName, ConfigFileName)
	if err != nil {
		return err
	}

	// Save config
	return shared.SaveConfig(config, configPath)
}

// GetDataDir returns the data directory for the client
func GetDataDir() (string, error) {
	configDir, err := shared.GetConfigDir(AppName)
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(configDir, "data")
	return dataDir, nil
}