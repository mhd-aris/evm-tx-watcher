package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Address     string     `json:"address" db:"address"`
	ChainID     int        `json:"chain_id" db:"chain_id"`
	IsContract  bool       `json:"is_contract" db:"is_contract"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	Label       *string    `json:"label,omitempty" db:"label"`
	Description *string    `json:"description,omitempty" db:"description"`
	UserID      *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
