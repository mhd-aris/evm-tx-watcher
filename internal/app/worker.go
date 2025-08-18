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

func RunWorker(ctx context.Context, cfg *config.Config, log *util.Logger) error {
	var wg sync.WaitGroup
	blockChan := make(chan *types.Block, 50) // change back to *types.Block

	// spawn watcher per network
	for name, rpcURL := range cfg.Networks {
		c, err := client.New(name, rpcURL, log)
		if err != nil {
			return err
		}

		// watcher now returns *types.Block directly
		w := watcher.New(c, 1, log)

		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			if err := w.Start(ctx, blockChan); err != nil && ctx.Err() == nil {
				log.Errorf("Watcher %s stopped: %v", n, err)
			}
		}(name)
	}

	// processor
	proc := processor.New(log)

	// consumer
	go func() {
		for block := range blockChan {
			if err := proc.HandleBlock(ctx, block); err != nil {
				log.Errorf("Processor error: %v", err)
			}
		}
	}()

	// shutdown gracefully
	go func() {
		wg.Wait()
		close(blockChan)
	}()

	<-ctx.Done()
	return nil
}
