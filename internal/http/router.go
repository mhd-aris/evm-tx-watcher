package http

import (
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/http/handler"
	"evm-tx-watcher/internal/http/middleware"
	"evm-tx-watcher/internal/repository"
	"evm-tx-watcher/internal/service"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(cfg *config.Config, logger *logrus.Logger, db *sqlx.DB) *echo.Echo {
	e := echo.New()

	e.HideBanner = false

	// Add middlewares
	e.Use(middleware.LoggingMiddleware(logger))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Setup routes
	setupRoutes(e, db)

	return e
}

func setupRoutes(e *echo.Echo, db *sqlx.DB) {

	// Health check endpoint
	e.GET("/health", handler.HealthHandler)

	addrRepo := repository.NewAddressRepository(db)
	addrService := service.NewAddressService(addrRepo)
	addrHandler := handler.NewAddressHandler(addrService)

	v1 := e.Group("/api/v1")
	{
		v1.GET("/addresses", addrHandler.GetAll)
		v1.POST("/addresses", addrHandler.Register)
	}
}
