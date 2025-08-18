package watcher

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"evm-tx-watcher/internal/blockchain/client"
	"evm-tx-watcher/internal/util"
)

type Watcher struct {
	client        *client.Client
	confirmations int64
	log           *util.Logger
}

func New(c *client.Client, conf int64, log *util.Logger) *Watcher {
	return &Watcher{
		client:        c,
		confirmations: conf,
		log:           log,
	}
}

func (w *Watcher) Start(ctx context.Context, out chan<- *types.Block) error {
	// sanity check: get latest head
	head, err := w.client.Eth.BlockNumber(ctx)
	if err != nil {
		w.log.WithError(err).Errorf("[%s] failed to get latest block number", w.client.Network)
	} else {
		w.log.Infof("[%s] starting watcher, current head=%d", w.client.Network, head)
	}

	headers := make(chan *types.Header)
	sub, err := w.client.Eth.SubscribeNewHead(ctx, headers)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			w.log.Infof("[%s] watcher stopped", w.client.Network)
			return nil

		case err := <-sub.Err():
			w.log.WithError(err).Warnf("[%s] subscription error, retrying in 5s", w.client.Network)
			select {
			case <-time.After(5 * time.Second):
				return w.Start(ctx, out) // restart watcher
			case <-ctx.Done():
				return nil
			}

		case header := <-headers:
			if header == nil {
				continue
			}

			// cek confirmations
			if w.confirmations > 0 {
				head, err := w.client.Eth.BlockNumber(ctx)
				if err != nil {
					w.log.WithError(err).Warnf("[%s] failed to get head", w.client.Network)
					continue
				}
				if head < header.Number.Uint64()+uint64(w.confirmations) {
					w.log.Debugf("[%s] skip block %d waiting for %d confirmations (head=%d)",
						w.client.Network,
						header.Number.Uint64(),
						w.confirmations,
						head,
					)
					continue
				}
			}

			// Try ethclient first, fallback to raw RPC for L2 networks
			block, err := w.client.Eth.BlockByNumber(ctx, header.Number)
			if err != nil {
				if err.Error() == "transaction type not supported" {
					// L2 networks like Arbitrum/Base use custom transaction types
					// Fall back to header-only information
					w.log.Infof("[%s] new block=%d hash=%s (L2 block, transaction details unavailable)",
						w.client.Network,
						header.Number.Uint64(),
						header.Hash().Hex(),
					)
					continue
				}
				w.log.WithError(err).Warnf("[%s] failed to fetch block %d", w.client.Network, header.Number.Uint64())
				continue
			}

			w.log.Infof("[%s] new block=%d hash=%s txs=%d",
				w.client.Network,
				header.Number.Uint64(),
				header.Hash().Hex(),
				len(block.Transactions()),
			)

			select {
			case out <- block: // send block
			default:
				w.log.Warnf("[%s] channel full, dropping block %d",
					w.client.Network, header.Number.Uint64())
			}
		}
	}
}
