package main

import (
	_ "evm-tx-watcher/docs"
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/db"
	"evm-tx-watcher/internal/http"
	"evm-tx-watcher/internal/util"
	"evm-tx-watcher/internal/validator"
	"fmt"
	"log"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           EVM Tx Watcher API
// @version         1.0
// @description     A simple REST API service for monitoring Ethereum wallet addresses and getting notified when transactions occur.
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger with config
	logger := util.NewLogger(cfg.LogLevel, cfg.LogFormat)
	logger.Info("Starting EVM Transaction Watcher server...")

	// Initialize request validator
	v := validator.NewValidator()

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
	e := http.NewRouter(cfg, db, logger, v)
	e.Validator = v

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	logger.Infof("Starting server on %s", addr)
	if err := e.Start(addr); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
