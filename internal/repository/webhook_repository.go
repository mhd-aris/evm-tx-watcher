package repository

import (
	"context"
	"database/sql"
	"errors"
	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WebhookRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, webhook *domain.Webhook) (*domain.Webhook, error)
	FindByID(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*domain.Webhook, error)
	FindByAddress(ctx context.Context, tx *sqlx.Tx, address string) (*domain.Webhook, error)
	FindAll(ctx context.Context, tx *sqlx.Tx) ([]*domain.Webhook, error)
}

type webhookRepository struct {
	db *sqlx.DB
}

func NewWebhookRepository(db *sqlx.DB) WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) Create(ctx context.Context, tx *sqlx.Tx, webhook *domain.Webhook) (*domain.Webhook, error) {

	_, err := tx.ExecContext(ctx, `INSERT INTO webhooks (id, address_id, url, secret, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, webhook.ID, webhook.AddressID, webhook.URL, webhook.Secret, webhook.CreatedAt, webhook.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return webhook, nil
}

func (r *webhookRepository) FindByID(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*domain.Webhook, error) {
	var webhook domain.Webhook
	query := `SELECT id, address, created_at, updated_at FROM webhooks WHERE id = $1`
	err := tx.GetContext(ctx, &webhook, query, id)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *webhookRepository) FindByAddress(ctx context.Context, tx *sqlx.Tx, address string) (*domain.Webhook, error) {
	var webhook domain.Webhook
	query := `SELECT address_id, address, created_at, updated_at FROM webhooks WHERE address = $1`
	err := tx.GetContext(ctx, &webhook, query, address)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *webhookRepository) FindAll(ctx context.Context, tx *sqlx.Tx) ([]*domain.Webhook, error) {
	var webhooks []*domain.Webhook
	query := `SELECT id, address, created_at, updated_at FROM webhooks`
	err := tx.SelectContext(ctx, &webhooks, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return webhooks, nil
}
