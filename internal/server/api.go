package server

import (
	"fmt"
	"net/http"
	"networkmonitor/shared"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// API handles the HTTP and WebSocket API
type API struct {
	clientManager *ClientManager
	router        *gin.Engine
	upgrader      websocket.Upgrader
}

// NewAPI creates a new API handler
func NewAPI(clientManager *ClientManager) *API {
	// Create router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create API
	api := &API{
		clientManager: clientManager,
		router:        router,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins
			},
		},
	}

	// Set up routes
	api.setupRoutes()

	return api
}

// Start starts the API server
func (a *API) Start(address string) error {
	return a.router.Run(address)
}

// setupRoutes sets up the API routes
func (a *API) setupRoutes() {
	// WebSocket endpoint
	a.router.GET("/ws", a.handleWebSocket)

	// Client API
	a.router.GET("/api/clients", a.getClients)
	a.router.GET("/api/clients/:id", a.getClient)
	a.router.GET("/api/clients/:id/requests", a.getClientRequests)

	// Config API
	a.router.GET("/api/config", a.getConfig)
	a.router.PUT("/api/config", a.updateConfig)

	// Static files
	a.router.Static("/dashboard", "./web/dist")
	a.router.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}

// handleWebSocket handles WebSocket connections
func (a *API) handleWebSocket(c *gin.Context) {
	// Upgrade connection to WebSocket
	ws, err := a.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Create client connection
	conn := NewClientConnection(ws, a.clientManager)

	// Start handling connection
	conn.Start()
}

// getClients returns a list of all clients
func (a *API) getClients(c *gin.Context) {
	clients := a.clientManager.GetClients()
	c.JSON(http.StatusOK, clients)
}

// getClient returns information about a specific client
func (a *API) getClient(c *gin.Context) {
	id := c.Param("id")
	client, found := a.clientManager.GetClient(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

// getClientRequests returns network requests for a specific client
func (a *API) getClientRequests(c *gin.Context) {
	id := c.Param("id")
	
	// Get query parameters
	limit := 100
	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			limit = 100
		}
	}
	
	// Get requests from storage
	requests, err := a.clientManager.storage.GetNetworkRequests(id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get requests"})
		return
	}
	
	c.JSON(http.StatusOK, requests)
}

// getConfig returns the server configuration
func (a *API) getConfig(c *gin.Context) {
	config, err := a.clientManager.storage.GetServerConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get config"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// updateConfig updates the server configuration
func (a *API) updateConfig(c *gin.Context) {
	var config shared.ServerConfig
	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config format"})
		return
	}
	
	if err := a.clientManager.storage.SaveServerConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}