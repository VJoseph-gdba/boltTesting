package server

import (
	"fmt"
	"networkmonitor/shared"
	"os"
	"path/filepath"
)

// Server is the main server application
type Server struct {
	config       shared.ServerConfig
	storage      *Storage
	clientManager *ClientManager
	api          *API
	systray      *SystrayHandler
}

// NewServer creates a new server instance
func NewServer() (*Server, error) {
	// Get data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	dataDir := filepath.Join(homeDir, ".config", "NetworkMonitor", "server")

	// Create storage
	storage, err := NewStorage(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	// Load config
	config, err := storage.GetServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create client manager
	clientManager := NewClientManager(storage)

	// Create API
	api := NewAPI(clientManager)

	// Create server
	server := &Server{
		config:       config,
		storage:      storage,
		clientManager: clientManager,
		api:          api,
	}

	// Create systray handler
	server.systray = NewSystrayHandler(server)

	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	fmt.Println("Starting Network Monitor Server...")

	// Start systray
	s.systray.Start()

	// Start API server
	fmt.Printf("Starting API server on %s\n", s.config.ListenAddress)
	return s.api.Start(s.config.ListenAddress)
}

// Stop stops the server
func (s *Server) Stop() {
	fmt.Println("Stopping Network Monitor Server...")

	// Stop systray
	s.systray.Stop()
}