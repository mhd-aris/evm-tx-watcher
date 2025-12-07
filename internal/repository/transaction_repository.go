package repository

import (
	"context"
	"database/sql"
	"fmt"

	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, transaction *domain.Transaction) (domain.Transaction, error)
	FindByHash(ctx context.Context, hash string) (*domain.Transaction, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	FindByBlockNumber(ctx context.Context, chainID int64, blockNumber int64) ([]*domain.Transaction, error)
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *sqlx.Tx, transaction *domain.Transaction) (domain.Transaction, error) {
	query := `
		INSERT INTO transactions (
			id, hash, block_number, block_hash, transaction_index, chain_id,
			from_address, to_address, value, gas_used, gas_price, tx_type,
			status, block_timestamp, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	_, err := tx.ExecContext(ctx, query,
		transaction.ID,
		transaction.Hash,
		transaction.BlockNumber,
		transaction.BlockHash,
		transaction.TransactionIndex,
		transaction.ChainID,
		transaction.FromAddress,
		transaction.ToAddress,
		transaction.Value.String(), // Store as string to handle big numbers
		transaction.GasUsed,
		transaction.GasPrice.String(), // Store as string to handle big numbers
		transaction.TxType,
		transaction.Status,
		transaction.BlockTimestamp,
		transaction.CreatedAt,
	)

	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to insert transaction: %w", err)
	}

	return *transaction, nil
}

func (r *transactionRepository) FindByHash(ctx context.Context, hash string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT id, hash, block_number, block_hash, transaction_index, chain_id,
		       from_address, to_address, value, gas_used, gas_price, tx_type,
		       status, block_timestamp, created_at
		FROM transactions 
		WHERE hash = $1`

	err := r.db.GetContext(ctx, &transaction, query, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find transaction by hash: %w", err)
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT id, hash, block_number, block_hash, transaction_index, chain_id,
		       from_address, to_address, value, gas_used, gas_price, tx_type,
		       status, block_timestamp, created_at
		FROM transactions 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find transaction by ID: %w", err)
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByBlockNumber(ctx context.Context, chainID int64, blockNumber int64) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	query := `
		SELECT id, hash, block_number, block_hash, transaction_index, chain_id,
		       from_address, to_address, value, gas_used, gas_price, tx_type,
		       status, block_timestamp, created_at
		FROM transactions 
		WHERE chain_id = $1 AND block_number = $2
		ORDER BY transaction_index`

	err := r.db.SelectContext(ctx, &transactions, query, chainID, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions by block number: %w", err)
	}

	return transactions, nil
}

