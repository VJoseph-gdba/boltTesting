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
# Clone the repository
git clone https://github.com/yourusername/network-monitor.git
cd network-monitor

# Build everything
make

# Or build components separately
make build-client
make build-server
make build-web
```

### Running the Server

```bash
make run-server
# or
./bin/server
```

### Running the Client

```bash
make run-client
# or
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



Base prompt:
create me an app in go that will have a client and a server the client will send https request to two website and send the complete result to the server (handshake, dnc, tcp etc...) the server will recieve these info and will process it by a list of client connected and disconnected. When we click on a client we have his dashboard where we can see since when he is online and what are the status of the client https request. The goal here is to have a graph and list to navigate throught the history of the client and to be able to notice and analyse network problem on the client side. Additionnaly the client and server will be manage through various command by using the systray. I also want the client to be able to modify config file from a software from distance, open it on the server webpage and edit it live without asking for administrator permission. i want the interface to be clean and modern like material design 3
