package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Logging middleware
func LoggingMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Execute next handler
			err := next(c)

			logEntry := logger.WithFields(logrus.Fields{
				"method": c.Request().Method,
				"path":   c.Request().URL.Path,
				"status": c.Response().Status,
			})

			if err != nil {
				logEntry.WithField("error", err.Error()).Error("Request failed")
			} else {
				logEntry.Info("Request completed")
			}

			return err
		}
	}
}
