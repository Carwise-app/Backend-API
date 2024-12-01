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

func init() {
	validate.RegisterValidation("strong_password", strongPassword)
	validate.RegisterValidation("password_match", validatePasswordMatch)
}

func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 || len(password) > 48 {
		return false
	}
	hasLower := false
	hasUpper := false
	hasDigit := false

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	return hasLower && hasUpper && hasDigit
}

func validatePasswordMatch(fl validator.FieldLevel) bool {
	password := fl.Parent().FieldByName("Password").String()
	rePassword := fl.Field().String()

	return password == rePassword
}
