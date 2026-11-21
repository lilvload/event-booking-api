package routes

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func validationMessage(err error) string {
	ve, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return "Could not parse request data."
	}

	fe := ve[0]

	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required."
	case "email":
		return "Invalid email format."
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters long."
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters long."
	default:
		return "Invalid value for " + fe.Field() + "."
	}
}
