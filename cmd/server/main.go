package main

import (
	"fmt"
	"networkmonitor/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create server
	server, err := server.NewServer()
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Stop server
	server.Stop()
}