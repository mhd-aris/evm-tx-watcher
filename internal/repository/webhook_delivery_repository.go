package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WebhookDeliveryRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, delivery *domain.WebhookDelivery) (domain.WebhookDelivery, error)
	Update(ctx context.Context, tx *sqlx.Tx, delivery *domain.WebhookDelivery) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.WebhookDelivery, error)
	FindPendingRetries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error)
	FindByWebhookID(ctx context.Context, webhookID uuid.UUID, limit int) ([]*domain.WebhookDelivery, error)
}

type webhookDeliveryRepository struct {
	db *sqlx.DB
}

func NewWebhookDeliveryRepository(db *sqlx.DB) WebhookDeliveryRepository {
	return &webhookDeliveryRepository{db: db}
}

func (r *webhookDeliveryRepository) Create(ctx context.Context, tx *sqlx.Tx, delivery *domain.WebhookDelivery) (domain.WebhookDelivery, error) {
	query := `
		INSERT INTO webhook_deliveries (
			id, webhook_id, transaction_id, payload, status, http_status_code,
			response_body, error_message, retry_count, max_retries, next_retry_at,
			delivered_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`

	_, err := tx.ExecContext(ctx, query,
		delivery.ID,
		delivery.WebhookID,
		delivery.TransactionID,
		delivery.Payload,
		delivery.Status,
		delivery.HTTPStatusCode,
		delivery.ResponseBody,
		delivery.ErrorMessage,
		delivery.RetryCount,
		delivery.MaxRetries,
		delivery.NextRetryAt,
		delivery.DeliveredAt,
		delivery.CreatedAt,
		delivery.UpdatedAt,
	)

	if err != nil {
		return domain.WebhookDelivery{}, fmt.Errorf("failed to insert webhook delivery: %w", err)
	}

	return *delivery, nil
}

func (r *webhookDeliveryRepository) Update(ctx context.Context, tx *sqlx.Tx, delivery *domain.WebhookDelivery) error {
	query := `
		UPDATE webhook_deliveries SET
			status = $2,
			http_status_code = $3,
			response_body = $4,
			error_message = $5,
			retry_count = $6,
			next_retry_at = $7,
			delivered_at = $8,
			updated_at = $9
		WHERE id = $1`

	delivery.UpdatedAt = time.Now()

	_, err := tx.ExecContext(ctx, query,
		delivery.ID,
		delivery.Status,
		delivery.HTTPStatusCode,
		delivery.ResponseBody,
		delivery.ErrorMessage,
		delivery.RetryCount,
		delivery.NextRetryAt,
		delivery.DeliveredAt,
		delivery.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update webhook delivery: %w", err)
	}

	return nil
}

func (r *webhookDeliveryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.WebhookDelivery, error) {
	var delivery domain.WebhookDelivery
	query := `
		SELECT id, webhook_id, transaction_id, payload, status, http_status_code,
		       response_body, error_message, retry_count, max_retries, next_retry_at,
		       delivered_at, created_at, updated_at
		FROM webhook_deliveries 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &delivery, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find webhook delivery by ID: %w", err)
	}

	return &delivery, nil
}

func (r *webhookDeliveryRepository) FindPendingRetries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error) {
	var deliveries []*domain.WebhookDelivery
	query := `
		SELECT id, webhook_id, transaction_id, payload, status, http_status_code,
		       response_body, error_message, retry_count, max_retries, next_retry_at,
		       delivered_at, created_at, updated_at
		FROM webhook_deliveries 
		WHERE status = 'failed' 
		  AND retry_count < max_retries 
		  AND (next_retry_at IS NULL OR next_retry_at <= NOW())
		ORDER BY created_at ASC
		LIMIT $1`

	err := r.db.SelectContext(ctx, &deliveries, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending retries: %w", err)
	}

	return deliveries, nil
}

func (r *webhookDeliveryRepository) FindByWebhookID(ctx context.Context, webhookID uuid.UUID, limit int) ([]*domain.WebhookDelivery, error) {
	var deliveries []*domain.WebhookDelivery
	query := `
		SELECT id, webhook_id, transaction_id, payload, status, http_status_code,
		       response_body, error_message, retry_count, max_retries, next_retry_at,
		       delivered_at, created_at, updated_at
		FROM webhook_deliveries 
		WHERE webhook_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	err := r.db.SelectContext(ctx, &deliveries, query, webhookID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find webhook deliveries by webhook ID: %w", err)
	}

	return deliveries, nil
}

