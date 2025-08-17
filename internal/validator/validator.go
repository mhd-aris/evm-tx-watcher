package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("eth_addr", validateAddress)

	return &Validator{
		validator: v,
	}
}

func (v *Validator) Validate(s interface{}) error {
	return v.validator.Struct(s)
}

func (v *Validator) ParseValidationError(err error) map[string]string {
	validationErrors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrs {
			field := strings.ToLower(validationErr.Field())
			tag := validationErr.Tag()

			// Custom error messages berdasarkan tag
			switch tag {
			case "required":
				validationErrors[field] = field + " is required"
			case "email":
				validationErrors[field] = field + " must be a valid email address"
			case "url":
				validationErrors[field] = field + " must be a valid URL"
			case "min":
				validationErrors[field] = field + " must be at least " + validationErr.Param() + " characters"
			case "max":
				validationErrors[field] = field + " must be at most " + validationErr.Param() + " characters"
			case "eth_addr":
				validationErrors[field] = field + " must be a valid Ethereum address"
			default:
				validationErrors[field] = field + " is invalid"
			}
		}
	}

	return validationErrors
}
