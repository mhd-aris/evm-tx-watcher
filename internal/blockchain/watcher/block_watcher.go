package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"evm-tx-watcher/internal/blockchain/client"
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/util"
)

// BlockEvent represents a confirmed block with its transactions
type BlockEvent struct {
	NetworkConfig      config.NetworkConfig
	Block              *types.Block
	TransactionDetails []*client.TransactionDetails
}

type Watcher struct {
	client        *client.Client
	confirmations int64
	logger        *util.Logger
	networkConfig config.NetworkConfig
}

func New(c *client.Client, networkConfig config.NetworkConfig, confirmations int64, logger *util.Logger) *Watcher {
	return &Watcher{
		client:        c,
		confirmations: confirmations,
		logger:        logger,
		networkConfig: networkConfig,
	}
}

func (w *Watcher) Start(ctx context.Context, out chan<- *types.Block) error {
	// Get latest block number for initial sync
	latestBlock, err := w.client.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block number for %s: %w", w.networkConfig.Name, err)
	}

	w.logger.Infof("[%s] Starting watcher, current head=%d, confirmations=%d",
		w.networkConfig.Name, latestBlock, w.confirmations)

	// Subscribe to new headers
	headers := make(chan *types.Header, 10)
	errCh, err := w.client.SubscribeNewHeads(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads for %s: %w", w.networkConfig.Name, err)
	}

	// Track pending blocks
	pending := make(map[uint64]*types.Header)

	for {
		select {
		case <-ctx.Done():
			w.logger.Infof("[%s] Watcher stopped", w.networkConfig.Name)
			return nil

		case err := <-errCh:
			if err != nil {
				w.logger.WithError(err).Warnf("[%s] Subscription error, retrying in 10s", w.networkConfig.Name)
				select {
				case <-time.After(10 * time.Second):
					return w.Start(ctx, out) // restart watcher
				case <-ctx.Done():
					return nil
				}
			}

		case header := <-headers:
			if header == nil {
				continue
			}

			w.logger.Debugf("[%s] New header: block=%d", w.networkConfig.Name, header.Number.Uint64())

			// Store new block header in pending
			pending[header.Number.Uint64()] = header
			currentHead := header.Number.Uint64()

			// Process confirmed blocks
			for blockNum, blockHeader := range pending {
				if currentHead >= blockNum+uint64(w.confirmations) {
					// Block is now confirmed
					if err := w.processConfirmedBlock(ctx, blockHeader, out); err != nil {
						w.logger.WithError(err).Errorf("[%s] Failed to process confirmed block %d",
							w.networkConfig.Name, blockNum)
					}
					delete(pending, blockNum)
				}
			}

			// Clean up old pending blocks (older than 100 blocks)
			for blockNum := range pending {
				if currentHead > blockNum+100 {
					w.logger.Warnf("[%s] Dropping old pending block %d", w.networkConfig.Name, blockNum)
					delete(pending, blockNum)
				}
			}
		}
	}
}

func (w *Watcher) processConfirmedBlock(ctx context.Context, header *types.Header, out chan<- *types.Block) error {
	// Get basic block via client method
	block, _, err := w.client.GetBlockWithTransactions(ctx, header.Number)
	if err != nil {
		return fmt.Errorf("failed to get block %d: %w", header.Number.Uint64(), err)
	}

	w.logger.Infof("[%s] Confirmed block=%d hash=%s txs=%d",
		w.networkConfig.Name, block.Number().Uint64(), block.Hash().Hex(), len(block.Transactions()))

	// Send to processor (non-blocking)
	select {
	case out <- block:
		return nil
	default:
		w.logger.Warnf("[%s] Block processor channel full, dropping block %d",
			w.networkConfig.Name, block.Number().Uint64())
		return fmt.Errorf("block processor channel full")
	}
}
