package client

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
)

// SystrayHandler manages system tray integration
type SystrayHandler struct {
	client   *Client
	quitChan chan struct{}
}

// NewSystrayHandler creates a new system tray handler
func NewSystrayHandler(client *Client) *SystrayHandler {
	return &SystrayHandler{
		client:   client,
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
	systray.SetTitle("Network Monitor Client")
	systray.SetTooltip("Network Monitor Client")

	// Set icon (this would be a real icon in a production app)
	// systray.SetIcon(iconData)

	// Create menu items
	mStatus := systray.AddMenuItem("Status: Initializing...", "Shows connection status")
	mStatus.Disable()
	systray.AddSeparator()

	mToggleConnection := systray.AddMenuItem("Disconnect", "Connect/Disconnect from server")
	mViewDashboard := systray.AddMenuItem("Open Dashboard", "Open the web dashboard")
	mSettings := systray.AddMenuItem("Settings", "Open settings")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Update status based on connection state
	go func() {
		for {
			if s.client.connection.IsConnected() {
				mStatus.SetTitle("Status: Connected")
				mToggleConnection.SetTitle("Disconnect")
			} else {
				mStatus.SetTitle("Status: Disconnected")
				mToggleConnection.SetTitle("Connect")
			}

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
			case <-mToggleConnection.ClickedCh:
				if s.client.connection.IsConnected() {
					s.client.connection.Disconnect()
				} else {
					s.client.connection.Connect()
				}

			case <-mViewDashboard.ClickedCh:
				// Open browser to dashboard
				url := fmt.Sprintf("%s/dashboard", s.client.config.ServerAddress)
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
	s.client.Stop()
	os.Exit(0)
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) {
	// Implementation depends on platform
	// This is just a placeholder
	fmt.Printf("Opening browser to: %s\n", url)
}