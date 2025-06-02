package shared

// ServerConfig represents the server configuration
type ServerConfig struct {
	MaxClients      int    `json:"maxClients"`
	HistoryDays     int    `json:"historyDays"`
	ListenAddress   string `json:"listenAddress"`
	RefreshInterval int    `json:"refreshInterval"`
}