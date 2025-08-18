package processor

import (
	"context"

	"evm-tx-watcher/internal/util"

	"github.com/ethereum/go-ethereum/core/types"
)

type Processor struct {
	log *util.Logger
}

func New(log *util.Logger) *Processor {
	return &Processor{log: log}
}

func (p *Processor) HandleBlock(ctx context.Context, blk *types.Block) error {
	p.log.Infof("[Processor] block=%d hash=%s txs=%d",
		blk.NumberU64(), blk.Hash().Hex(), len(blk.Transactions()))

	// Process transactions if available
	if len(blk.Transactions()) > 0 {
		for i, tx := range blk.Transactions() {
			// Limit logging to first 3 transactions to avoid spam
			if i >= 3 {
				p.log.Infof("[Processor] ... and %d more transactions", len(blk.Transactions())-3)
				break
			}

			p.log.Infof("[Processor] tx[%d]: hash=%s from=%s to=%s value=%s gas=%d type=%d",
				i,
				tx.Hash().Hex(),
				p.getSender(tx),
				p.getRecipient(tx),
				tx.Value().String(),
				tx.Gas(),
				tx.Type(),
			)
		}
	} else {
		p.log.Infof("[Processor] no transaction details available for this block")
	}

	return nil
}

func (p *Processor) getSender(tx *types.Transaction) string {
	// For now, return placeholder since we'd need chain ID and signature recovery
	// In a full implementation, you'd use types.Sender() with appropriate signer
	return "0x..." // placeholder
}

func (p *Processor) getRecipient(tx *types.Transaction) string {
	if tx.To() != nil {
		return tx.To().Hex()
	}
	return "contract_creation"
}
