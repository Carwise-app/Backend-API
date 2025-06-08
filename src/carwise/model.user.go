package carwise

type User struct {
	Id          string `json:"id"`
	GoogleId    string `json:"google_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	ImageUrl    string `json:"image_url"`
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	Role        int    `json:"role"`
	Status      int    `json:"status"`
	DeviceToken string `json:"-"`
	EmailNotify bool   `json:"-"`
	PushNotify  bool   `json:"-"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	LastLogin   int64  `json:"last_login"`
}
