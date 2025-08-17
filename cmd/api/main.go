package main

import (
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/db"
	"evm-tx-watcher/internal/http"
	"evm-tx-watcher/internal/util"
	"fmt"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           EVM Tx Watcher API
// @version         1.0
// @description     API untuk register address dan webhook listener.
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger with config
	logger := util.NewLogger(cfg.LogLevel)
	logger.Info("Starting EVM Transaction Watcher server...")

	// Initialize database connection
	db, err := db.InitDB(&cfg.DB)
	if err != nil {
		logger.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Create server address
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Infof("Server will listen on %s", addr)

	// Initialize Echo router
	e := http.NewRouter(cfg, logger, db)

	// Swagger documentation
	e.File("/swagger/doc.json", "docs/swagger.json")
	e.Static("/swagger", "docs")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	logger.Infof("Starting server on %s", addr)
	if err := e.Start(addr); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
