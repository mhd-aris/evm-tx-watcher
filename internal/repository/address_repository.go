package repository

import (
	"context"
	"database/sql"
	"errors"

	"evm-tx-watcher/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AddressRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, address *domain.Address) (domain.Address, error)
	// Update(ctx context.Context, tx *sqlx.Tx, address *domain.Address) error
	// Delete(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	FindByAddress(ctx context.Context, address string) (*domain.Address, error)
	FindAll(ctx context.Context) ([]*domain.Address, error)
}

type addressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, tx *sqlx.Tx, address *domain.Address) (domain.Address, error) {

	_, err := tx.ExecContext(ctx, `
    INSERT INTO addresses (id, address, chain_id, is_contract, is_active, label, description, user_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		address.ID,
		address.Address,
		address.ChainID,
		address.IsContract,
		address.IsActive,
		address.Label,
		address.Description,
		address.UserID,
		address.CreatedAt,
		address.UpdatedAt,
	)
	if err != nil {
		return domain.Address{}, err
	}

	return *address, nil
}

func (r *addressRepository) FindAll(ctx context.Context) ([]*domain.Address, error) {
	var addresses []*domain.Address
	query := `SELECT id, address, chain_id, is_contract, is_active, created_at, updated_at, label, description, user_id FROM addresses`
	err := r.db.SelectContext(ctx, &addresses, query)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	query := `SELECT id, address, chain_id, is_contract, is_active, created_at, updated_at, label, description, user_id FROM addresses WHERE id = $1`
	err := r.db.GetContext(ctx, &address, query, id)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) FindByAddress(ctx context.Context, addr string) (*domain.Address, error) {
	var address domain.Address
	query := `SELECT id, address, chain_id, is_contract, is_active, created_at, updated_at, label, description, user_id FROM addresses WHERE address = $1`
	err := r.db.GetContext(ctx, &address, query, addr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}
