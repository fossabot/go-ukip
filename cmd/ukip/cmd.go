package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanjay7178/go-ukip/internal/config"
	"github.com/sanjay7178/go-ukip/internal/device"
	"github.com/sanjay7178/go-ukip/internal/logging"
)

func main() {
	// Parse command-line flags
	_ = flag.String("config", "/etc/ukip/config.json", "Path to configuration file")
	flag.Parse()

	// Initialize logging
	if err := logging.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
		os.Exit(1)
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
	if err := monitor.Start(); err != nil {
		logging.Log.Fatalf("Failed to start device monitor: %v", err)
	}

	logging.Log.Info("UKIP started successfully")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan

	logging.Log.Info("Shutting down UKIP...")

	// Stop the monitor
	monitor.Stop()

	logging.Log.Info("UKIP stopped")
}