// Package carwise provides the data transfer objects and business logic for the Carwise application
package carwise

// @model UserCreateRequest
// @Description User registration request
type UserCreateRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50" example:"John"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50" example:"Doe"`
	CountryCode string `json:"country_code" validate:"required,max=10" example:"+90"`
	PhoneNumber string `json:"phone_number" validate:"required" example:"5551234567"`
	Email       string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password    string `json:"password" validate:"required,strong_password" example:"StrongP@ss123"`
}

// @model UserLoginRequest
// @Description User login request
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"StrongP@ss123"`
}

// @model ResetPasswordRequest
// @Description Password reset request
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

// @model ChangePasswordRequest
// @Description Change password request
type ChangePasswordRequest struct {
	Password   string `json:"password" validate:"required,strong_password" example:"NewStrongP@ss123"`
	RePassword string `json:"re_password" validate:"required,strong_password,password_match" example:"NewStrongP@ss123"`
}

// @model UserInfo
// @Description User information
type UserInfo struct {
	Id          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName   string `json:"first_name" example:"John"`
	LastName    string `json:"last_name" example:"Doe"`
	Email       string `json:"email" example:"john.doe@example.com"`
	CountryCode string `json:"country_code" example:"+90"`
	PhoneNumber string `json:"phone_number" example:"5551234567"`
}

// @model TokenResponse
// @Description JWT token response
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
