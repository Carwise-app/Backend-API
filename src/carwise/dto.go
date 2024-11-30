package carwise

type UserCreateRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	CountryCode string `json:"country_code" validate:"required,max=10"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type BrandResponse struct {
	Id     int              `json:"id"`
	Name   string           `json:"name"`
	Count  int              `json:"count"`
	Series []SeriesResponse `json:"series"`
}

type ModelResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}
type SeriesResponse struct {
	Id     int             `json:"id"`
	Name   string          `json:"name"`
	Count  int             `json:"count"`
	Models []ModelResponse `json:"models"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"  validate:"required,email"`
}

type ChangePasswordRequest struct {
	Password   string `json:"password" validate:"required"`
	RePassword string `json:"re_password" validate:"required"`
}
