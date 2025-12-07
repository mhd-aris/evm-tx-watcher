package repository

import (
	"context"
	"fmt"

	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TokenTransferRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, transfer *domain.TokenTransfer) (domain.TokenTransfer, error)
	FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*domain.TokenTransfer, error)
	FindByTokenAddress(ctx context.Context, tokenAddress string, chainID int64) ([]*domain.TokenTransfer, error)
}

type tokenTransferRepository struct {
	db *sqlx.DB
}

func NewTokenTransferRepository(db *sqlx.DB) TokenTransferRepository {
	return &tokenTransferRepository{db: db}
}

func (r *tokenTransferRepository) Create(ctx context.Context, tx *sqlx.Tx, transfer *domain.TokenTransfer) (domain.TokenTransfer, error) {
	query := `
		INSERT INTO token_transfers (
			id, transaction_id, log_index, token_address, from_address, to_address,
			value, token_decimals, token_symbol, token_name, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)`

	_, err := tx.ExecContext(ctx, query,
		transfer.ID,
		transfer.TransactionID,
		transfer.LogIndex,
		transfer.TokenAddress,
		transfer.FromAddress,
		transfer.ToAddress,
		transfer.Value.String(), // Store as string to handle big numbers
		transfer.TokenDecimals,
		transfer.TokenSymbol,
		transfer.TokenName,
		transfer.CreatedAt,
	)

	if err != nil {
		return domain.TokenTransfer{}, fmt.Errorf("failed to insert token transfer: %w", err)
	}

	return *transfer, nil
}

func (r *tokenTransferRepository) FindByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*domain.TokenTransfer, error) {
	var transfers []*domain.TokenTransfer
	query := `
		SELECT id, transaction_id, log_index, token_address, from_address, to_address,
		       value, token_decimals, token_symbol, token_name, created_at
		FROM token_transfers 
		WHERE transaction_id = $1
		ORDER BY log_index`

	err := r.db.SelectContext(ctx, &transfers, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find token transfers by transaction ID: %w", err)
	}

	return transfers, nil
}

func (r *tokenTransferRepository) FindByTokenAddress(ctx context.Context, tokenAddress string, chainID int64) ([]*domain.TokenTransfer, error) {
	var transfers []*domain.TokenTransfer
	query := `
		SELECT tt.id, tt.transaction_id, tt.log_index, tt.token_address, tt.from_address, tt.to_address,
		       tt.value, tt.token_decimals, tt.token_symbol, tt.token_name, tt.created_at
		FROM token_transfers tt
		INNER JOIN transactions t ON tt.transaction_id = t.id
		WHERE tt.token_address = $1 AND t.chain_id = $2
		ORDER BY t.block_timestamp DESC, tt.log_index`

	err := r.db.SelectContext(ctx, &transfers, query, tokenAddress, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to find token transfers by token address: %w", err)
	}

	return transfers, nil
}
