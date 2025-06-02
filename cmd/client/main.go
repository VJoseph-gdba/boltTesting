package main

import (
	"fmt"
	"networkmonitor/internal/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create client
	client, err := client.NewClient()
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	// Start client
	if err := client.Start(); err != nil {
		fmt.Printf("Error starting client: %v\n", err)
		os.Exit(1)
	}

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Stop client
	client.Stop()
}