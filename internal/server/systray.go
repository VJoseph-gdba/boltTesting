package server

import (
	"fmt"
	"networkmonitor/shared"
	"os"

	"github.com/getlantern/systray"
)

// SystrayHandler manages system tray integration
type SystrayHandler struct {
	server   *Server
	quitChan chan struct{}
}

// NewSystrayHandler creates a new system tray handler
func NewSystrayHandler(server *Server) *SystrayHandler {
	return &SystrayHandler{
		server:   server,
		quitChan: make(chan struct{}),
	}
}

// Start initializes the system tray
func (s *SystrayHandler) Start() {
	go systray.Run(s.onReady, s.onExit)
}

// Stop stops the system tray
func (s *SystrayHandler) Stop() {
	close(s.quitChan)
	systray.Quit()
}

// onReady is called when the system tray is ready
func (s *SystrayHandler) onReady() {
	systray.SetTitle("Network Monitor Server")
	systray.SetTooltip("Network Monitor Server")

	// Set icon (this would be a real icon in a production app)
	// systray.SetIcon(iconData)

	// Create menu items
	mStatus := systray.AddMenuItem("Server Running", "Shows server status")
	mStatus.Disable()
	systray.AddSeparator()

	mClients := systray.AddMenuItem("Connected Clients: 0", "Shows connected clients")
	mClients.Disable()
	mOpenDashboard := systray.AddMenuItem("Open Dashboard", "Open the web dashboard")
	mSettings := systray.AddMenuItem("Settings", "Open settings")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Update connected clients count
	go func() {
		for {
			// Update clients count
			clients := s.server.clientManager.GetClients()
			connectedCount := 0
			for _, client := range clients {
				if client.Status == shared.StatusOnline {
					connectedCount++
				}
			}
			mClients.SetTitle(fmt.Sprintf("Connected Clients: %d", connectedCount))

			// Check for quit signal
			select {
			case <-s.quitChan:
				return
			default:
				// Continue looping
			}
		}
	}()

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mOpenDashboard.ClickedCh:
				// Open browser to dashboard
				url := fmt.Sprintf("http://%s/dashboard", s.server.config.ListenAddress)
				openBrowser(url)

			case <-mSettings.ClickedCh:
				// Open settings dialog
				fmt.Println("Settings clicked")

			case <-mQuit.ClickedCh:
				systray.Quit()
				return

			case <-s.quitChan:
				return
			}
		}
	}()
}

// onExit is called when the system tray is exiting
func (s *SystrayHandler) onExit() {
	// Cleanup
	s.server.Stop()
	os.Exit(0)
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) {
	// Implementation depends on platform
	// This is just a placeholder
	fmt.Printf("Opening browser to: %s\n", url)
}