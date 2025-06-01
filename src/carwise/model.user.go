package carwise

type User struct {
	Id          string
	GoogleId    string
	FirstName   string
	LastName    string
	ImageUrl    string
	CountryCode string
	PhoneNumber string
	Email       string
	Password    string
	Role        int
	Status      int
	DeviceToken string
	EmailNotify bool
	PushNotify  bool
	CreatedAt   int64
	UpdatedAt   int64
	LastLogin   int64
}
