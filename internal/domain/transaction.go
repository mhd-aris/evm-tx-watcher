package domain

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	Hash             string       `json:"hash" db:"hash"`
	BlockNumber      int64        `json:"block_number" db:"block_number"`
	BlockHash        string       `json:"block_hash" db:"block_hash"`
	TransactionIndex int          `json:"transaction_index" db:"transaction_index"`
	ChainID          int64        `json:"chain_id" db:"chain_id"`
	FromAddress      string       `json:"from_address" db:"from_address"`
	ToAddress        *string      `json:"to_address,omitempty" db:"to_address"`
	Value            *big.Int     `json:"value" db:"value"` // Wei amount for ETH transfers
	GasUsed          *int64       `json:"gas_used,omitempty" db:"gas_used"`
	GasPrice         *big.Int     `json:"gas_price,omitempty" db:"gas_price"`
	TxType           int          `json:"tx_type" db:"tx_type"`
	Status           int          `json:"status" db:"status"` // 1=success, 0=failed
	BlockTimestamp   time.Time    `json:"block_timestamp" db:"block_timestamp"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	TokenTransfers   []TokenTransfer `json:"token_transfers,omitempty" db:"-"`
}

// TokenTransfer represents an ERC-20 token transfer within a transaction
type TokenTransfer struct {
	ID            uuid.UUID `json:"id" db:"id"`
	TransactionID uuid.UUID `json:"transaction_id" db:"transaction_id"`
	LogIndex      int       `json:"log_index" db:"log_index"`
	TokenAddress  string    `json:"token_address" db:"token_address"`
	FromAddress   string    `json:"from_address" db:"from_address"`
	ToAddress     string    `json:"to_address" db:"to_address"`
	Value         *big.Int  `json:"value" db:"value"` // Raw token amount
	TokenDecimals *int      `json:"token_decimals,omitempty" db:"token_decimals"`
	TokenSymbol   *string   `json:"token_symbol,omitempty" db:"token_symbol"`
	TokenName     *string   `json:"token_name,omitempty" db:"token_name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID             uuid.UUID `json:"id" db:"id"`
	WebhookID      uuid.UUID `json:"webhook_id" db:"webhook_id"`
	TransactionID  uuid.UUID `json:"transaction_id" db:"transaction_id"`
	Payload        string    `json:"payload" db:"payload"` // JSON string
	Status         string    `json:"status" db:"status"`
	HTTPStatusCode *int      `json:"http_status_code,omitempty" db:"http_status_code"`
	ResponseBody   *string   `json:"response_body,omitempty" db:"response_body"`
	ErrorMessage   *string   `json:"error_message,omitempty" db:"error_message"`
	RetryCount     int       `json:"retry_count" db:"retry_count"`
	MaxRetries     int       `json:"max_retries" db:"max_retries"`
	NextRetryAt    *time.Time `json:"next_retry_at,omitempty" db:"next_retry_at"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// WebhookDeliveryStatus represents the status of webhook delivery
type WebhookDeliveryStatus string

const (
	WebhookDeliveryStatusPending            WebhookDeliveryStatus = "pending"
	WebhookDeliveryStatusDelivered          WebhookDeliveryStatus = "delivered"
	WebhookDeliveryStatusFailed             WebhookDeliveryStatus = "failed"
	WebhookDeliveryStatusMaxRetriesExceeded WebhookDeliveryStatus = "max_retries_exceeded"
)

// WatchedAddress represents an address being monitored
type WatchedAddress struct {
	Address    string `json:"address"`
	ChainID    int64  `json:"chain_id"`
	IsActive   bool   `json:"is_active"`
	WebhookID  uuid.UUID `json:"webhook_id"`
	WebhookURL string `json:"webhook_url"`
}

