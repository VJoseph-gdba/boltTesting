package server

import (
	"encoding/json"
	"fmt"
	"networkmonitor/shared"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientConnection represents a WebSocket connection to a client
type ClientConnection struct {
	ws           *websocket.Conn
	clientID     string
	clientInfo   shared.ClientInfo
	clientMgr    *ClientManager
	sendChan     chan shared.ServerMessage
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewClientConnection creates a new client connection
func NewClientConnection(ws *websocket.Conn, clientMgr *ClientManager) *ClientConnection {
	return &ClientConnection{
		ws:        ws,
		clientMgr: clientMgr,
		sendChan:  make(chan shared.ServerMessage, 100),
		stopChan:  make(chan struct{}),
	}
}

// Start begins handling the client connection
func (c *ClientConnection) Start() {
	// Start goroutines for sending and receiving
	c.wg.Add(2)
	go c.sendLoop()
	go c.receiveLoop()
}

// Stop stops handling the client connection
func (c *ClientConnection) Stop() {
	close(c.stopChan)
	c.ws.Close()
	c.wg.Wait()
}

// SendMessage sends a message to the client
func (c *ClientConnection) SendMessage(msgType string, data interface{}) {
	msg := shared.ServerMessage{
		Type:      msgType,
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

// sendLoop sends messages to the client
func (c *ClientConnection) sendLoop() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			return
		case msg := <-c.sendChan:
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error marshaling message: %v\n", err)
				continue
			}

			if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				return
			}
		}
	}
}

// receiveLoop receives messages from the client
func (c *ClientConnection) receiveLoop() {
	defer c.wg.Done()
	defer c.clientMgr.RemoveClient(c.clientID)

	for {
		select {
		case <-c.stopChan:
			return
		default:
			// Set read deadline
			c.ws.SetReadDeadline(time.Now().Add(time.Minute * 2))

			// Read message
			_, message, err := c.ws.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading message: %v\n", err)
				return
			}

			// Process message
			var clientMsg shared.ClientMessage
			if err := json.Unmarshal(message, &clientMsg); err != nil {
				fmt.Printf("Error unmarshaling message: %v\n", err)
				continue
			}

			// Set client ID if not set
			if c.clientID == "" {
				c.clientID = clientMsg.ClientID
			}

			// Handle message based on type
			c.handleClientMessage(clientMsg)
		}
	}
}

// handleClientMessage processes messages from the client
func (c *ClientConnection) handleClientMessage(msg shared.ClientMessage) {
	switch msg.Type {
	case shared.TypeClientConnect:
		// Handle client connect
		if clientInfo, ok := msg.Data.(map[string]interface{}); ok {
			// Convert to ClientInfo
			infoBytes, _ := json.Marshal(clientInfo)
			json.Unmarshal(infoBytes, &c.clientInfo)
			
			c.clientInfo.Status = shared.StatusOnline
			c.clientInfo.LastSeen = time.Now()
			c.clientMgr.AddClient(c.clientID, c)
			
			// Store client info
			c.clientMgr.storage.SaveClientInfo(c.clientInfo)
		}

	case shared.TypeHeartbeat:
		// Handle heartbeat
		if clientInfo, ok := msg.Data.(map[string]interface{}); ok {
			// Convert to ClientInfo
			infoBytes, _ := json.Marshal(clientInfo)
			json.Unmarshal(infoBytes, &c.clientInfo)
			
			c.clientInfo.Status = shared.StatusOnline
			c.clientInfo.LastSeen = time.Now()
			
			// Update client info
			c.clientMgr.storage.SaveClientInfo(c.clientInfo)
		}

	case shared.TypeClientDisconnect:
		// Handle client disconnect
		if clientInfo, ok := msg.Data.(map[string]interface{}); ok {
			// Convert to ClientInfo
			infoBytes, _ := json.Marshal(clientInfo)
			json.Unmarshal(infoBytes, &c.clientInfo)
			
			c.clientInfo.Status = shared.StatusOffline
			
			// Update client info
			c.clientMgr.storage.SaveClientInfo(c.clientInfo)
			c.clientMgr.RemoveClient(c.clientID)
		}

	case shared.TypeNetworkRequest:
		// Handle network request
		if requestData, ok := msg.Data.(map[string]interface{}); ok {
			// Convert to NetworkRequest
			var request shared.NetworkRequest
			requestBytes, _ := json.Marshal(requestData)
			json.Unmarshal(requestBytes, &request)
			
			// Store request
			c.clientMgr.storage.SaveNetworkRequest(c.clientID, request)
		}

	default:
		fmt.Printf("Unknown message type: %s\n", msg.Type)
	}
}

// ClientManager manages client connections
type ClientManager struct {
	clients     map[string]*ClientConnection
	storage     *Storage
	mutex       sync.RWMutex
}

// NewClientManager creates a new client manager
func NewClientManager(storage *Storage) *ClientManager {
	return &ClientManager{
		clients: make(map[string]*ClientConnection),
		storage: storage,
	}
}

// AddClient adds a client connection
func (m *ClientManager) AddClient(clientID string, conn *ClientConnection) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Remove existing connection if any
	if existing, found := m.clients[clientID]; found {
		existing.Stop()
	}

	// Add new connection
	m.clients[clientID] = conn

	fmt.Printf("Client connected: %s\n", clientID)
}

// RemoveClient removes a client connection
func (m *ClientManager) RemoveClient(clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, found := m.clients[clientID]; found {
		delete(m.clients, clientID)
		fmt.Printf("Client disconnected: %s\n", clientID)
	}
}

// GetClient gets a client by ID
func (m *ClientManager) GetClient(clientID string) (shared.ClientInfo, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if conn, found := m.clients[clientID]; found {
		return conn.clientInfo, true
	}

	// Check storage for offline clients
	client, err := m.storage.GetClientInfo(clientID)
	if err != nil {
		return shared.ClientInfo{}, false
	}

	return client, true
}

// GetClients gets all clients
func (m *ClientManager) GetClients() []shared.ClientInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Get online clients
	clients := make([]shared.ClientInfo, 0, len(m.clients))
	for _, conn := range m.clients {
		clients = append(clients, conn.clientInfo)
	}

	// Get offline clients from storage
	offlineClients, err := m.storage.GetAllClientInfo()
	if err != nil {
		return clients
	}

	// Create a map of online clients for quick lookup
	onlineMap := make(map[string]bool)
	for _, client := range clients {
		onlineMap[client.ID] = true
	}

	// Add offline clients that aren't already in the list
	for _, client := range offlineClients {
		if !onlineMap[client.ID] {
			clients = append(clients, client)
		}
	}

	return clients
}

// SendMessageToClient sends a message to a specific client
func (m *ClientManager) SendMessageToClient(clientID, msgType string, data interface{}) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if conn, found := m.clients[clientID]; found {
		conn.SendMessage(msgType, data)
		return true
	}

	return false
}

// BroadcastMessage sends a message to all connected clients
func (m *ClientManager) BroadcastMessage(msgType string, data interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, conn := range m.clients {
		conn.SendMessage(msgType, data)
	}
}