package http

import (
	"evm-tx-watcher/internal/config"
	"evm-tx-watcher/internal/http/handler"
	"evm-tx-watcher/internal/http/middleware"
	"evm-tx-watcher/internal/repository"
	"evm-tx-watcher/internal/service"
	"evm-tx-watcher/internal/util"
	"evm-tx-watcher/internal/validator"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func NewRouter(cfg *config.Config, db *sqlx.DB, logger *util.Logger, validator *validator.Validator) *echo.Echo {
	e := echo.New()

	e.HideBanner = false

	// Add middlewares
	e.Use(middleware.LoggingMiddleware(logger))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Setup routes
	setupRoutes(e, db, logger, validator)

	return e
}

func setupRoutes(e *echo.Echo, db *sqlx.DB, logger *util.Logger, validator *validator.Validator) {

	// Health check endpoint
	e.GET("/health", handler.HealthHandler)

	unitOfWork := repository.NewUnitOfWork(db)
	addrRepo := repository.NewAddressRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	addrService := service.NewAddressService(unitOfWork, addrRepo, webhookRepo)
	addrHandler := handler.NewAddressHandler(addrService, logger, validator)

	v1 := e.Group("/api/v1")
	{
		v1.GET("/addresses", addrHandler.GetAll)
		v1.POST("/addresses", addrHandler.Register)
	}
}
