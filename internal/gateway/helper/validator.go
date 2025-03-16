package helper

import "github.com/go-playground/validator/v10"

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}
