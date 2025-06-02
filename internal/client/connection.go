package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"networkmonitor/shared"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Connection manages the WebSocket connection to the server
type Connection struct {
	serverURL     string
	ws            *websocket.Conn
	clientInfo    shared.ClientInfo
	sendChan      chan shared.ClientMessage
	stopChan      chan struct{}
	reconnectChan chan struct{}
	wg            sync.WaitGroup
	connected     bool
	mutex         sync.Mutex
}

// NewConnection creates a new server connection
func NewConnection(serverAddress, clientName string) *Connection {
	hostname, _ := os.Hostname()
	if clientName == "" {
		clientName = hostname
	}

	// Generate a stable client ID based on hostname
	clientID := uuid.NewMD5(uuid.NameSpaceDNS, []byte(hostname)).String()

	return &Connection{
		serverURL: serverAddress,
		clientInfo: shared.ClientInfo{
			ID:          clientID,
			Name:        clientName,
			Status:      shared.StatusOffline,
			ConnectedAt: time.Time{},
			LastSeen:    time.Time{},
			Version:     "1.0.0",
			OSInfo:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		},
		sendChan:      make(chan shared.ClientMessage, 100),
		stopChan:      make(chan struct{}),
		reconnectChan: make(chan struct{}, 1),
	}
}

// Connect establishes a connection to the server
func (c *Connection) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connected {
		return nil
	}

	// Parse server URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return err
	}

	// Convert to WebSocket URL
	wsURL := url.URL{Scheme: "ws", Host: u.Host, Path: "/ws"}
	if u.Scheme == "https" {
		wsURL.Scheme = "wss"
	}

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		return err
	}

	c.ws = ws
	c.connected = true
	c.clientInfo.Status = shared.StatusOnline
	c.clientInfo.ConnectedAt = time.Now()
	c.clientInfo.LastSeen = time.Now()

	// Start goroutines for sending and receiving
	c.wg.Add(2)
	go c.sendLoop()
	go c.receiveLoop()

	// Send connect message
	c.SendMessage(shared.TypeClientConnect, c.clientInfo)

	return nil
}

// Disconnect closes the connection to the server
func (c *Connection) Disconnect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.connected {
		return
	}

	// Send disconnect message
	disconnectTime := time.Now()
	c.clientInfo.DisconnectedAt = &disconnectTime
	c.clientInfo.Status = shared.StatusOffline
	c.SendMessage(shared.TypeClientDisconnect, c.clientInfo)

	// Signal stop to goroutines
	close(c.stopChan)

	// Close WebSocket
	c.ws.Close()
	c.connected = false

	// Wait for goroutines to finish
	c.wg.Wait()

	// Reset channels
	c.stopChan = make(chan struct{})
	c.sendChan = make(chan shared.ClientMessage, 100)
}

// IsConnected returns whether the client is connected
func (c *Connection) IsConnected() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.connected
}

// SendMessage sends a message to the server
func (c *Connection) SendMessage(msgType string, data interface{}) {
	msg := shared.ClientMessage{
		Type:      msgType,
		ClientID:  c.clientInfo.ID,
		Timestamp: time.Now(),
		Data:      data,
	}

	select {
	case c.sendChan <- msg:
		// Message queued successfully
	default:
		// Channel full, log error
		fmt.Printf("Warning: send channel full, dropping message of type %s\n", msgType)
	}
}

// sendLoop sends messages to the server
func (c *Connection) sendLoop() {
	defer c.wg.Done()

	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case msg := <-c.sendChan:
			c.mutex.Lock()
			if !c.connected {
				c.mutex.Unlock()
				continue
			}

			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error marshaling message: %v\n", err)
				c.mutex.Unlock()
				continue
			}

			err = c.ws.WriteMessage(websocket.TextMessage, data)
			c.mutex.Unlock()

			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				c.triggerReconnect()
				return
			}

		case <-heartbeatTicker.C:
			c.mutex.Lock()
			if !c.connected {
				c.mutex.Unlock()
				continue
			}

			// Update last seen time
			c.clientInfo.LastSeen = time.Now()

			// Send heartbeat
			heartbeat := shared.ClientMessage{
				Type:      shared.TypeHeartbeat,
				ClientID:  c.clientInfo.ID,
				Timestamp: time.Now(),
				Data:      c.clientInfo,
			}

			data, err := json.Marshal(heartbeat)
			if err != nil {
				fmt.Printf("Error marshaling heartbeat: %v\n", err)
				c.mutex.Unlock()
				continue
			}

			err = c.ws.WriteMessage(websocket.TextMessage, data)
			c.mutex.Unlock()

			if err != nil {
				fmt.Printf("Error sending heartbeat: %v\n", err)
				c.triggerReconnect()
				return
			}
		}
	}
}

// receiveLoop receives messages from the server
func (c *Connection) receiveLoop() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			return
		default:
			c.mutex.Lock()
			if !c.connected {
				c.mutex.Unlock()
				time.Sleep(time.Second)
				continue
			}

			// Set read deadline
			c.ws.SetReadDeadline(time.Now().Add(time.Minute))
			c.mutex.Unlock()

			// Read message
			_, message, err := c.ws.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading message: %v\n", err)
				c.triggerReconnect()
				return
			}

			// Process message
			var serverMsg shared.ServerMessage
			if err := json.Unmarshal(message, &serverMsg); err != nil {
				fmt.Printf("Error unmarshaling message: %v\n", err)
				continue
			}

			// Handle message based on type
			c.handleServerMessage(serverMsg)
		}
	}
}

// handleServerMessage processes messages from the server
func (c *Connection) handleServerMessage(msg shared.ServerMessage) {
	switch msg.Type {
	case shared.TypeConfigRequest:
		// Handle config request
		fmt.Println("Config request received from server")
		// Implementation would go here

	case shared.TypeCommandRequest:
		// Handle command request
		fmt.Println("Command request received from server")
		// Implementation would go here

	default:
		fmt.Printf("Unknown message type: %s\n", msg.Type)
	}
}

// triggerReconnect schedules a reconnection attempt
func (c *Connection) triggerReconnect() {
	c.mutex.Lock()
	c.connected = false
	c.mutex.Unlock()

	select {
	case c.reconnectChan <- struct{}{}:
		// Reconnect signal sent
	default:
		// Channel already has reconnect pending
	}
}

// ManageConnection handles connection management including reconnection
func (c *Connection) ManageConnection() {
	// Initial connection
	if err := c.Connect(); err != nil {
		fmt.Printf("Initial connection failed: %v\n", err)
		c.triggerReconnect()
	}

	// Reconnection loop
	for {
		select {
		case <-c.reconnectChan:
			// Wait before reconnecting
			time.Sleep(5 * time.Second)

			// Try to disconnect if still connected
			c.Disconnect()

			// Try to reconnect
			backoff := time.Second
			maxBackoff := time.Minute
			for i := 0; i < 10; i++ {
				fmt.Printf("Attempting to reconnect (%d/10)...\n", i+1)
				if err := c.Connect(); err == nil {
					fmt.Println("Reconnected successfully")
					break
				} else {
					fmt.Printf("Reconnection failed: %v\n", err)
					time.Sleep(backoff)
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				}
			}
		}
	}
}