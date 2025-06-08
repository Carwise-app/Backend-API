// Package carwise provides the data transfer objects and business logic for the Carwise application
package carwise

// @model ProfileResponse
// @Description User profile response
type ProfileResponse struct {
	Id          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FirstName   string `json:"first_name" example:"John"`
	LastName    string `json:"last_name" example:"Doe"`
	ImageUrl    string `json:"image_url" example:"https://example.com/avatar.jpg"`
	CountryCode string `json:"country_code" example:"+90"`
	PhoneNumber string `json:"phone_number" example:"5551234567"`
	Email       string `json:"email" example:"john.doe@example.com"`
	Role        int    `json:"role" example:"1"`
	Status      int    `json:"status" example:"1"`
	CreatedAt   int64  `json:"created_at" example:"1646092800"`
}

// @model ProfileEditRequest
// @Description Profile edit request
type ProfileEditRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50" example:"John"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50" example:"Doe"`
	CountryCode string `json:"country_code" validate:"required,max=10" example:"+90"`
	PhoneNumber string `json:"phone_number" validate:"required" example:"5551234567"`
}

// @model ProfileNotifyRequest
// @Description Profile notify request
type ProfileNotifyRequest struct {
	DeviceToken string `json:"device_token" example:"123e4567-e89b-12d3-a456-426614174000"`
	EmailNotify bool   `json:"email_notify" example:"true"`
	PushNotify  bool   `json:"push_notify" example:"true"`
	UserId      string `json:"-"`
}

// @model GetProfileNotifyRequest
// @Description Get profile notify Request
type GetProfileNotifyRequest struct {
	UserId string `json:"-"`
}

// @model GetProfileNotifyResponse
// @Description Get profile notify response
type GetProfileNotifyResponse struct {
	DeviceToken string `json:"device_token" example:"123e4567-e89b-12d3-a456-426614174000"`
	EmailNotify bool   `json:"email_notify" example:"true"`
	PushNotify  bool   `json:"push_notify" example:"true"`
}

// @model DeleteAccountRequest
// @Description Delete account request
type DeleteAccountRequest struct {
	ProfileId string `json:"-"`
	UserId    string `json:"-"`
	Role      int    `json:"-"`
}
