package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version,omitempty"`
}

func HealthHandler(c echo.Context) error {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "evm-tx-watcher",
		Version:   "1.0.0",
	}

	return c.JSON(http.StatusOK, response)
}
