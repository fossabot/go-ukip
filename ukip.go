package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanjay7178/go-ukip/device"
	"github.com/sanjay7178/go-ukip/logging"
	"github.com/sanjay7178/go-ukip/config"
)

func main() {
	// Initialize logging
	if err := logging.Init(); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logging.Log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize device monitor
	monitor, err := device.NewMonitor(cfg)
	if err != nil {
		logging.Log.Fatalf("Failed to initialize device monitor: %v", err)
	}

	// Start monitoring
	go monitor.Start()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Cleanup and exit
	monitor.Stop()
	logging.Log.Info("UKIP stopped")
}
