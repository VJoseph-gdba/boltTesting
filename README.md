# Network Monitor

A comprehensive Go application for monitoring network requests and analyzing client network health.

## Features

- **Client-Server Architecture**: Real-time monitoring of multiple clients
- **HTTPS Request Monitoring**: Track detailed metrics for HTTP requests
- **Network Diagnostics**: Analyze DNS, TCP, TLS handshake, and other connection metrics
- **Interactive Dashboard**: View client status and network performance data
- **Remote Configuration**: Edit client configuration files through the web interface
- **System Tray Integration**: Manage client and server applications through the system tray

## Requirements

- Go 1.18 or later
- Node.js 14 or later (for the web interface)
- Modern web browser

## Quick Start

### Building the Application

```bash
# Build the client
go build -o bin/client ./cmd/client

# Build the server
go build -o bin/server ./cmd/server

# Build the web interface
cd web && npm install && npm run build
```

### Running the Server

```bash
./bin/server
```

### Running the Client

```bash
./bin/client
```

## Client Features

- Sends HTTPS requests to configured websites
- Collects detailed networking information (handshake, DNS, TCP)
- Sends metrics to the server
- System tray interface for management
- Supports remote configuration

## Server Features

- Receives and processes client data
- Tracks connected/disconnected clients
- Provides dashboards for each client
- Web interface with Material Design 3 UI
- System tray integration

## Architecture

```
├── cmd/                  # Entry points
│   ├── client/           # Client application
│   └── server/           # Server application
├── internal/             # Internal packages
│   ├── client/           # Client implementation
│   └── server/           # Server implementation
├── shared/               # Shared types and utilities
├── web/                  # Web interface
│   ├── src/              # React frontend
│   └── dist/             # Built frontend
└── bin/                  # Compiled binaries
```

## Configuration

### Client Configuration

The client configuration is stored in `~/.config/NetworkMonitor/client.json` with the following structure:

```json
{
  "serverAddress": "http://localhost:8080",
  "clientName": "MyClient",
  "targets": [
    {
      "name": "Google",
      "url": "https://www.google.com",
      "interval": 60,
      "enabled": true
    },
    {
      "name": "GitHub",
      "url": "https://github.com",
      "interval": 60,
      "enabled": true
    }
  ],
  "logLevel": "info"
}
```

### Server Configuration

The server configuration is stored in `~/.config/NetworkMonitor/server/config.json` with the following structure:

```json
{
  "maxClients": 100,
  "historyDays": 30,
  "listenAddress": ":8080",
  "refreshInterval": 60
}
```

## License

MIT