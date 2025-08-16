package http

import (
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/http/handler"
	"evm-tx-watcher/internal/http/middleware"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(cfg *config.Config, logger *logrus.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = false

	// Add middlewares
	e.Use(middleware.LoggingMiddleware(logger))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Setup routes
	setupRoutes(e)

	return e
}

func setupRoutes(e *echo.Echo) {

	// Health check endpoint
	e.GET("/health", handler.HealthHandler)

	v1 := e.Group("/api/v1")
	{
		_ = v1 // Avoid unused variable warning for now
	}
}
