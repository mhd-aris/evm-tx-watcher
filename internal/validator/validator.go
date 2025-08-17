package validator

import (
	v10 "github.com/go-playground/validator/v10"
)

var Validate *v10.Validate

func New() *v10.Validate {
	Validate = v10.New()

	Validate.RegisterValidation("eth_addr", validateAddress)
	return Validate
}
