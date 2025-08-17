package handler

import (
	"evm-tx-watcher/internal/dto"
	"evm-tx-watcher/internal/errors"
	"evm-tx-watcher/internal/http/response"
	"evm-tx-watcher/internal/service"
	"evm-tx-watcher/internal/util"
	"evm-tx-watcher/internal/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	addressService service.AddressService
	logger         *util.Logger
	validator      *validator.Validator
}

func NewAddressHandler(
	addressService service.AddressService,
	logger *util.Logger,
	validator *validator.Validator,
) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		logger:         logger,
		validator:      validator,
	}
}

// RegisterAddress godoc
// @Summary      Register address to monitor
// @Description  Adds an EVM address and webhook for transaction monitoring
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        payload body dto.RegisterAddressRequest true "Register Request"
// @Success      201 {object} dto.AddressResponse
// @Failure      400 {object} dto.BaseResponse
// @Router       /addresses [post]
func (h *AddressHandler) Register(c echo.Context) error {
	var request dto.RegisterAddressRequest

	if err := c.Bind(&request); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		appErr := errors.ValidationError("Invalid JSON format")
		return response.SendAppError(c, appErr)
	}

	if err := h.validator.Validate(request); err != nil {
		h.logger.WithError(err).Error("Failed to validate request")
		return response.SendValidationError(c, h.validator, err)

	}

	createdAddress, err := h.addressService.Register(c.Request().Context(), &request)
	if err != nil {
		h.logger.WithError(err).Error("Failed to register address")
		return response.SendAppError(c, err)
	}
	return response.SendSuccess(c, http.StatusCreated, "Address registered successfully", createdAddress)
}

// GetAll godoc
// @Summary      Get all registered addresses
// @Description  Retrieve a list of all registered addresses
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.AddressResponse
// @Failure      400 {object} dto.BaseResponse
// @Router       /addresses [get]
func (h *AddressHandler) GetAll(c echo.Context) error {
	addresses, err := h.addressService.GetAll(c.Request().Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve addresses")
		return response.SendAppError(c, err)
	}

	return response.SendSuccess(c, http.StatusOK, "Addresses retrieved successfully", addresses)
}
