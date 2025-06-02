package client

import (
	"fmt"
	"networkmonitor/shared"
	"sync"
)

// Client is the main client application
type Client struct {
	config     shared.ClientConfig
	monitor    *Monitor
	connection *Connection
	systray    *SystrayHandler
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewClient creates a new client instance
func NewClient() (*Client, error) {
	// Load configuration
	config, err := LoadClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create components
	monitor := NewMonitor()
	connection := NewConnection(config.ServerAddress, config.ClientName)

	// Create client
	client := &Client{
		config:     config,
		monitor:    monitor,
		connection: connection,
		stopChan:   make(chan struct{}),
	}

	// Create systray handler
	client.systray = NewSystrayHandler(client)

	return client, nil
}

// Start starts the client
func (c *Client) Start() error {
	fmt.Println("Starting Network Monitor Client...")

	// Start monitor
	c.monitor.Start(c.config.Targets)

	// Start connection management
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.connection.ManageConnection()
	}()

	// Start processing network requests
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.processNetworkRequests()
	}()

	// Start systray
	c.systray.Start()

	fmt.Println("Network Monitor Client started")
	return nil
}

// Stop stops the client
func (c *Client) Stop() {
	fmt.Println("Stopping Network Monitor Client...")

	// Stop systray
	c.systray.Stop()

	// Stop monitor
	c.monitor.Stop()

	// Disconnect from server
	c.connection.Disconnect()

	// Signal stop to goroutines
	close(c.stopChan)

	// Wait for goroutines to finish
	c.wg.Wait()

	fmt.Println("Network Monitor Client stopped")
}

// UpdateConfig updates the client configuration
func (c *Client) UpdateConfig(config shared.ClientConfig) error {
	// Save configuration
	if err := SaveClientConfig(config); err != nil {
		return err
	}

	// Update client configuration
	c.config = config

	// Restart monitor with new targets
	c.monitor.Start(config.Targets)

	// Update connection if server address changed
	if c.connection.serverURL != config.ServerAddress {
		c.connection.Disconnect()
		c.connection = NewConnection(config.ServerAddress, config.ClientName)
		c.connection.Connect()
	}

	return nil
}

// processNetworkRequests processes network requests from the monitor
func (c *Client) processNetworkRequests() {
	resultChan := c.monitor.GetResultChan()

	for {
		select {
		case <-c.stopChan:
			return
		case result := <-resultChan:
			// Send network request to server if connected
			if c.connection.IsConnected() {
				c.connection.SendMessage(shared.TypeNetworkRequest, result)
			}
			
			// Log result
			status := result.StatusCode
			if result.Error != "" {
				fmt.Printf("[%s] Error: %s - %s (%s)\n", 
					result.TargetName, result.URL, result.Error, result.ErrorType)
			} else {
				fmt.Printf("[%s] %d: %s (%dms)\n", 
					result.TargetName, status, result.URL, result.TotalTime)
			}
		}
	}
}