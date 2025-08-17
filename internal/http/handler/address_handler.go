package handler

import (
	"evm-tx-watcher/internal/dto"
	"evm-tx-watcher/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	service service.AddressService
}

func NewAddressHandler(service service.AddressService) *AddressHandler {
	return &AddressHandler{service: service}
}

// RegisterAddress godoc
// @Summary      Register address to monitor
// @Description  Adds an EVM address and webhook for transaction monitoring
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        payload body dto.RegisterRequest true "Register Request"
// @Success      201 {object} dto.RegisterResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /addresses [post]
func (h *AddressHandler) Register(c echo.Context) error {
	var request dto.RegisterRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request format",
		})
	}

	createdAddress, err := h.service.Register(c.Request().Context(), &request)
	if err != nil {
		// TODO: Add proper error type checking for different error types
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := dto.RegisterResponse{
		ID:      createdAddress.ID.String(),
		Address: createdAddress.Address,
		Label:   *createdAddress.Label,
	}

	return c.JSON(http.StatusCreated, response)
}

// GetAll godoc
// @Summary      Get all registered addresses
// @Description  Retrieve a list of all registered addresses
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.RegisterResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /addresses [get]
func (h *AddressHandler) GetAll(c echo.Context) error {
	addresses, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to retrieve addresses",
		})
	}

	// Convert to response DTOs
	var responses []dto.RegisterResponse
	for _, addr := range addresses {
		response := dto.RegisterResponse{
			ID:      addr.ID.String(),
			Address: addr.Address,
		}
		if addr.Label != nil {
			response.Label = *addr.Label
		}
		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, responses)
}
