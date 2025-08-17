package response

import (
	"evm-tx-watcher/internal/dto"
	"evm-tx-watcher/internal/errors"
	"evm-tx-watcher/internal/validator"

	"github.com/labstack/echo/v4"
)

func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	response := dto.SuccessResponse(message, data)
	return c.JSON(statusCode, response)
}

func SendError(c echo.Context, statusCode int, appErr *errors.AppError) error {
	response := dto.ErrorResponse(string(appErr.Code), appErr.Message, appErr.Details)
	return c.JSON(statusCode, response)
}

func SendValidationError(c echo.Context, validator *validator.Validator, err error) error {
	validationErrors := validator.ParseValidationError(err)
	response := dto.ValidationErrorResponse(validationErrors)
	return c.JSON(400, response)
}

func SendAppError(c echo.Context, appErr *errors.AppError) error {
	statusCode := GetStatusCodeFromError(appErr.Code)
	return SendError(c, statusCode, appErr)
}

func GetStatusCodeFromError(code errors.ErrorCode) int {
	switch code {
	case errors.ErrCodeValidation, errors.ErrCodeBadRequest:
		return 400
	case errors.ErrCodeNotFound:
		return 404
	case errors.ErrCodeAlreadyExists:
		return 409
	case errors.ErrCodeUnauthorized:
		return 401
	default:
		return 500
	}
}
