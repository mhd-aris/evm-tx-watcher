package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"evm-tx-watcher/internal/app"
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/util"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := util.NewLogger(cfg.LogLevel, cfg.LogFormat)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle signals
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		logger.Info("shutting down worker...")
		cancel()
	}()

	if err := app.RunWorker(ctx, cfg, logger); err != nil {
		logger.WithError(err).Error("worker exited with error")
		os.Exit(1)
	}
}
