package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) []string {
	err := validate.Struct(s)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed validation: %s\n", err.Field(), err.Tag())
			errorMessages = append(errorMessages, msg)
		}
		return errorMessages
	}
	return nil
}
