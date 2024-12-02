package main

import (
	"carwise"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) []string {
	err := validate.Struct(s)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag())
			errorMessages = append(errorMessages, msg)
		}
		return errorMessages
	}
	return nil
}

func init() {
	validate.RegisterValidation("strong_password", strongPassword)
	validate.RegisterValidation("password_match", validatePasswordMatch)
	validate.RegisterValidation("currency", validateCurrency)
	validate.RegisterValidation("fuel_type", validateFuelType)
	validate.RegisterValidation("transmission", validateTransmission)
	validate.RegisterValidation("body_type", validateBodyType)
	validate.RegisterValidation("condition", validateCondition)
	validate.RegisterValidation("drive_type", validateDriveType)
	validate.RegisterValidation("seller_type", validateSellerType)
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

func validateCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return currency == carwise.CurrencyTRY || currency == carwise.CurrencyUSD || currency == carwise.CurrencyEUR
}

func validateFuelType(fl validator.FieldLevel) bool {
	fuelType := fl.Field().String()
	return fuelType == carwise.FuelTypeDiesel || fuelType == carwise.FuelTypePetrol || fuelType == carwise.FuelTypePetrolAndLPG || fuelType == carwise.FuelTypeHybrid || fuelType == carwise.FuelTypeElectric
}

func validateBodyType(fl validator.FieldLevel) bool {
	bodyType := fl.Field().String()
	return bodyType == carwise.BodyTypeSedan || bodyType == carwise.BodyTypeHatchback3 || bodyType == carwise.BodyTypeHatchback5 || bodyType == carwise.BodyTypeCoupe || bodyType == carwise.BodyTypeCabrio || bodyType == carwise.BodyTypeMPV || bodyType == carwise.BodyTypePickup || bodyType == carwise.BodyTypeRoadster || bodyType == carwise.BodyTypeStationWagon || bodyType == carwise.BodyTypeSUV
}

func validateCondition(fl validator.FieldLevel) bool {
	condition := fl.Field().String()
	return condition == carwise.PartConditionOriginal || condition == carwise.PartConditionPainted || condition == carwise.PartConditionChanged
}

func validateDriveType(fl validator.FieldLevel) bool {
	driveType := fl.Field().String()
	return driveType == carwise.DriveTypeFrontWheelDrive || driveType == carwise.DriveTypeRearWheelDrive || driveType == carwise.DriveTypeFourWheelDrive || driveType == carwise.DriveTypeAllWheelDrive
}

func validateSellerType(fl validator.FieldLevel) bool {
	sellerType := fl.Field().String()
	return sellerType == carwise.SellerTypeIndividual || sellerType == carwise.SellerTypeDealer
}

func validateTransmission(fl validator.FieldLevel) bool {
	transmissionType := fl.Field().String()
	return transmissionType == carwise.TransmissionAutomatic || transmissionType == carwise.TransmissionManual || transmissionType == carwise.TransmissionSemiautomatic
}
