package repository

import (
	"context"
	"database/sql"
	"fmt"

	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WebhookRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, webhook *domain.Webhook) (domain.Webhook, error)
	Update(ctx context.Context, tx *sqlx.Tx, webhook *domain.Webhook) error
	Delete(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Webhook, error)
	FindByAddressID(ctx context.Context, addressID uuid.UUID) ([]*domain.Webhook, error)
}

type webhookRepository struct {
	db *sqlx.DB
}

func NewWebhookRepository(db *sqlx.DB) WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) Create(ctx context.Context, tx *sqlx.Tx, webhook *domain.Webhook) (domain.Webhook, error) {
	query := `
		INSERT INTO webhooks (id, address_id, url, secret, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query,
		webhook.ID,
		webhook.AddressID,
		webhook.URL,
		webhook.Secret,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	)

	if err != nil {
		return domain.Webhook{}, fmt.Errorf("failed to insert webhook: %w", err)
	}

	return *webhook, nil
}

func (r *webhookRepository) Update(ctx context.Context, tx *sqlx.Tx, webhook *domain.Webhook) error {
	query := `
		UPDATE webhooks SET
			url = $2,
			secret = $3,
			updated_at = $4
		WHERE id = $1`

	_, err := tx.ExecContext(ctx, query,
		webhook.ID,
		webhook.URL,
		webhook.Secret,
		webhook.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}

func (r *webhookRepository) Delete(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	query := `DELETE FROM webhooks WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

func (r *webhookRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Webhook, error) {
	var webhook domain.Webhook
	query := `
		SELECT id, address_id, url, secret, created_at, updated_at
		FROM webhooks 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &webhook, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find webhook by ID: %w", err)
	}

	return &webhook, nil
}

func (r *webhookRepository) FindByAddressID(ctx context.Context, addressID uuid.UUID) ([]*domain.Webhook, error) {
	var webhooks []*domain.Webhook
	query := `
		SELECT id, address_id, url, secret, created_at, updated_at
		FROM webhooks 
		WHERE address_id = $1
		ORDER BY created_at ASC`

	err := r.db.SelectContext(ctx, &webhooks, query, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to find webhooks by address ID: %w", err)
	}

	return webhooks, nil
}