package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validateAddress(fl validator.FieldLevel) bool {
	address := fl.Field().String()
	match, _ := regexp.MatchString("^0x[a-fA-F0-9]{40}$", address)
	return match
}
