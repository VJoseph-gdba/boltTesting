package shared

import (
	"time"
)

// ClientStatus represents the connection status of a client
type ClientStatus string

const (
	StatusOnline  ClientStatus = "online"
	StatusOffline ClientStatus = "offline"
)

// NetworkRequest represents a captured HTTP request
type NetworkRequest struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	Method        string    `json:"method"`
	StatusCode    int       `json:"statusCode"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	DNSTime       int64     `json:"dnsTime"`       // in milliseconds
	TCPTime       int64     `json:"tcpTime"`       // in milliseconds
	TLSTime       int64     `json:"tlsTime"`       // in milliseconds
	RequestTime   int64     `json:"requestTime"`   // in milliseconds
	ResponseTime  int64     `json:"responseTime"`  // in milliseconds
	TotalTime     int64     `json:"totalTime"`     // in milliseconds
	Error         string    `json:"error"`
	ErrorType     string    `json:"errorType"`
	TargetName    string    `json:"targetName"`    // Name of the monitored target
}

// ClientInfo represents information about a client
type ClientInfo struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	IPAddress    string       `json:"ipAddress"`
	Status       ClientStatus `json:"status"`
	ConnectedAt  time.Time    `json:"connectedAt"`
	DisconnectedAt *time.Time  `json:"disconnectedAt,omitempty"`
	LastSeen     time.Time    `json:"lastSeen"`
	Version      string       `json:"version"`
	OSInfo       string       `json:"osInfo"`
}

// ClientMessage represents a message sent from client to server
type ClientMessage struct {
	Type      string          `json:"type"`
	ClientID  string          `json:"clientId"`
	Timestamp time.Time       `json:"timestamp"`
	Data      interface{}     `json:"data"`
}

// ServerMessage represents a message sent from server to client
type ServerMessage struct {
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      interface{}     `json:"data"`
}

// ConfigFile represents a configuration file
type ConfigFile struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Content  string    `json:"content"`
	LastEdit time.Time `json:"lastEdit"`
}

// Target represents a website to monitor
type Target struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval int    `json:"interval"` // in seconds
	Enabled  bool   `json:"enabled"`
}

// ClientConfig represents the client configuration
type ClientConfig struct {
	ServerAddress string   `json:"serverAddress"`
	ClientName    string   `json:"clientName"`
	Targets       []Target `json:"targets"`
	LogLevel      string   `json:"logLevel"`
}

// MessageType constants
const (
	TypeHeartbeat          = "heartbeat"
	TypeNetworkRequest     = "network_request"
	TypeConfigUpdate       = "config_update"
	TypeConfigRequest      = "config_request"
	TypeConfigResponse     = "config_response"
	TypeCommandRequest     = "command_request"
	TypeCommandResponse    = "command_response"
	TypeClientConnect      = "client_connect"
	TypeClientDisconnect   = "client_disconnect"
	TypeClientsList        = "clients_list"
	TypeNetworkRequestList = "network_request_list"
)