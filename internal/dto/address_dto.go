package dto

import "time"

type RegisterAddressRequest struct {
	Address     string  `json:"address" validate:"required,eth_addr"`
	ChainID     int     `json:"chain_id" validate:"required,gt=0"`
	WebhookURL  string  `json:"webhook_url" validate:"required,url"`
	Secret      string  `json:"secret" validate:"required,min=10"`
	Label       *string `json:"label,omitempty" validate:"omitempty,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type AddressResponse struct {
	ID          string    `json:"id"`
	Address     string    `json:"address"`
	IsContract  bool      `json:"is_contract"`
	IsActive    bool      `json:"is_active"`
	Label       *string   `json:"label,omitempty"`
	Description *string   `json:"description,omitempty"`
	ChainID     int       `json:"chain_id"`
	WebhookURL  string    `json:"webhook_url"`
	UserID      *string   `json:"user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
