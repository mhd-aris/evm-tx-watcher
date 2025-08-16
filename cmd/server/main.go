package main

import (
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/http"
	"evm-tx-watcher/internal/util"
	"fmt"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger with config
	logger := util.NewLogger(cfg.LogLevel)
	logger.Info("Starting EVM Transaction Watcher server...")

	// Create server address
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Infof("Server will listen on %s", addr)

	// Initialize Echo router
	e := http.NewRouter(cfg, logger)

	// Start server
	logger.Infof("Starting server on %s", addr)
	if err := e.Start(addr); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
