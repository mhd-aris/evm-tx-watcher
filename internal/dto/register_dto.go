package dto

// RegisterRequest represents the request to register an address
type RegisterRequest struct {
	Address    string `json:"address" validate:"required,eth_addr"`
	ChainID    int    `json:"chain_id" validate:"required,gt=0"`
	WebhookURL string `json:"webhook_url" validate:"required,url"`
	Secret     string `json:"secret" validate:"required,min=10"`
	Label      string `json:"label,omitempty" validate:"omitempty,max=100"`
}

// RegisterResponse represents the response after registering an address
type RegisterResponse struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Label   string `json:"label,omitempty"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}
