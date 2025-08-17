package domain

import (
	"time"

	"github.com/google/uuid"
)

type Webhook struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AddressID uuid.UUID `json:"address_id" db:"address_id"`
	URL       string    `json:"url" db:"url"`
	Secret    string    `json:"-" db:"secret"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
