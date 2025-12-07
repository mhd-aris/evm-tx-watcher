package app

import (
	"context"
	"sync"

	"evm-tx-watcher/internal/blockchain/client"
	"evm-tx-watcher/internal/blockchain/watcher"
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/processor"
	"evm-tx-watcher/internal/util"

	"github.com/ethereum/go-ethereum/core/types"
)

func RunWorker(ctx context.Context, cfg *config.Config, logger *util.Logger) error {
	logger.Info("Starting EVM Transaction Watcher Worker")

	var wg sync.WaitGroup
	blockChan := make(chan *types.Block, 50)

	// Initialize simple processor
	proc := processor.New(logger)

	// Start blockchain watchers for each network
	for _, networkConfig := range cfg.Networks {
		logger.Infof("Initializing client for %s (Chain ID: %d)", networkConfig.Name, networkConfig.ChainID)

		// Create blockchain client
		blockchainClient, err := client.New(networkConfig, logger)
		if err != nil {
			logger.WithError(err).Errorf("Failed to create client for %s, skipping...", networkConfig.Name)
			continue // Skip this network, don't fail entire worker
		}

		// Create watcher with simple parameters
		blockWatcher := watcher.New(blockchainClient, networkConfig, 5, logger)

		// Start watcher in goroutine
		wg.Add(1)
		go func(network config.NetworkConfig, client *client.Client) {
			defer wg.Done()
			defer client.Close()

			logger.Infof("Starting watcher for %s", network.Name)

			if err := blockWatcher.Start(ctx, blockChan); err != nil && ctx.Err() == nil {
				logger.WithError(err).Errorf("Watcher %s stopped with error", network.Name)
			}
		}(networkConfig, blockchainClient)
	}

	// Start block processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Starting block processor")

		for {
			select {
			case <-ctx.Done():
				logger.Info("Block processor stopped")
				return
			case block := <-blockChan:
				if block != nil {
					if err := proc.HandleBlock(ctx, block); err != nil {
						logger.WithError(err).Error("Failed to process block")
					}
				}
			}
		}
	}()

	// Graceful shutdown
	go func() {
		wg.Wait()
		close(blockChan)
	}()

	// Wait for context cancellation
	<-ctx.Done()
	logger.Info("Shutting down worker...")

	return nil
}
