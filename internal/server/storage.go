package server

import (
	"encoding/json"
	"fmt"
	"networkmonitor/shared"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Storage handles data persistence
type Storage struct {
	dataDir      string
	clientsDir   string
	requestsDir  string
	configFile   string
	mutex        sync.RWMutex
}

// NewStorage creates a new storage handler
func NewStorage(dataDir string) (*Storage, error) {
	// Create directory structure
	clientsDir := filepath.Join(dataDir, "clients")
	requestsDir := filepath.Join(dataDir, "requests")
	configFile := filepath.Join(dataDir, "config.json")

	dirs := []string{dataDir, clientsDir, requestsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return &Storage{
		dataDir:     dataDir,
		clientsDir:  clientsDir,
		requestsDir: requestsDir,
		configFile:  configFile,
	}, nil
}

// SaveClientInfo saves client information
func (s *Storage) SaveClientInfo(client shared.ClientInfo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	filename := filepath.Join(s.clientsDir, client.ID+".json")
	data, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// GetClientInfo gets client information
func (s *Storage) GetClientInfo(clientID string) (shared.ClientInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	filename := filepath.Join(s.clientsDir, clientID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return shared.ClientInfo{}, err
	}

	var client shared.ClientInfo
	if err := json.Unmarshal(data, &client); err != nil {
		return shared.ClientInfo{}, err
	}

	return client, nil
}

// GetAllClientInfo gets all client information
func (s *Storage) GetAllClientInfo() ([]shared.ClientInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Read client directory
	files, err := os.ReadDir(s.clientsDir)
	if err != nil {
		return nil, err
	}

	// Read each client file
	clients := make([]shared.ClientInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filename := filepath.Join(s.clientsDir, file.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var client shared.ClientInfo
		if err := json.Unmarshal(data, &client); err != nil {
			continue
		}

		clients = append(clients, client)
	}

	return clients, nil
}

// SaveNetworkRequest saves a network request
func (s *Storage) SaveNetworkRequest(clientID string, request shared.NetworkRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create client requests directory
	clientDir := filepath.Join(s.requestsDir, clientID)
	if err := os.MkdirAll(clientDir, 0755); err != nil {
		return err
	}

	// Create date-based directory
	dateDir := filepath.Join(clientDir, request.StartTime.Format("2006-01-02"))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return err
	}

	// Create request file
	filename := filepath.Join(dateDir, request.ID+".json")
	data, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// GetNetworkRequests gets network requests for a client
func (s *Storage) GetNetworkRequests(clientID string, limit int) ([]shared.NetworkRequest, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Check if client directory exists
	clientDir := filepath.Join(s.requestsDir, clientID)
	if _, err := os.Stat(clientDir); os.IsNotExist(err) {
		return []shared.NetworkRequest{}, nil
	}

	// Read date directories
	dateDirs, err := os.ReadDir(clientDir)
	if err != nil {
		return nil, err
	}

	// Sort date directories in reverse order (newest first)
	// This is a simplification; in a real app, parse dates and sort properly
	var allRequests []shared.NetworkRequest

	// Process each date directory
	for _, dateDir := range dateDirs {
		if !dateDir.IsDir() {
			continue
		}

		// Read request files
		datePathFull := filepath.Join(clientDir, dateDir.Name())
		files, err := os.ReadDir(datePathFull)
		if err != nil {
			continue
		}

		// Process each request file
		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
				continue
			}

			// Read file
			filename := filepath.Join(datePathFull, file.Name())
			data, err := os.ReadFile(filename)
			if err != nil {
				continue
			}

			// Parse request
			var request shared.NetworkRequest
			if err := json.Unmarshal(data, &request); err != nil {
				continue
			}

			allRequests = append(allRequests, request)
			if len(allRequests) >= limit {
				break
			}
		}

		if len(allRequests) >= limit {
			break
		}
	}

	return allRequests, nil
}

// GetServerConfig gets the server configuration
func (s *Storage) GetServerConfig() (shared.ServerConfig, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Check if config file exists
	if _, err := os.Stat(s.configFile); os.IsNotExist(err) {
		// Create default config
		config := shared.ServerConfig{
			MaxClients:     100,
			HistoryDays:    30,
			ListenAddress:  ":8080",
			RefreshInterval: 60,
		}
		
		// Save default config
		if err := s.SaveServerConfig(config); err != nil {
			return config, err
		}
		
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(s.configFile)
	if err != nil {
		return shared.ServerConfig{}, err
	}

	// Parse config
	var config shared.ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return shared.ServerConfig{}, err
	}

	return config, nil
}

// SaveServerConfig saves the server configuration
func (s *Storage) SaveServerConfig(config shared.ServerConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(s.configFile, data, 0644)
}