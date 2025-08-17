package repository

import (
	"context"
	"time"

	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) (domain.Address, error)
	GetAll(ctx context.Context) ([]*domain.Address, error)
	// GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	// Update(ctx context.Context, address *domain.Address) error
	// Delete(ctx context.Context, id uuid.UUID) error
}

type addressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *domain.Address) (domain.Address, error) {
	address.ID = uuid.New()
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()
	address.IsActive = true // Default to active

	query := `INSERT INTO addresses (id, address, chain_id, is_contract, is_active, created_at, updated_at, label, description, user_id)
			  VALUES (:id, :address, :chain_id, :is_contract, :is_active, :created_at, :updated_at, :label, :description, :user_id)`
	_, err := r.db.NamedExecContext(ctx, query, address)
	if err != nil {
		return domain.Address{}, err
	}
	return *address, nil
}

func (r *addressRepository) GetAll(ctx context.Context) ([]*domain.Address, error) {
	var addresses []*domain.Address
	query := `SELECT id, address, chain_id, is_contract, is_active, created_at, updated_at, label, description, user_id FROM addresses`
	err := r.db.SelectContext(ctx, &addresses, query)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}
