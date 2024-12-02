package carwise

import "time"

type UserCreateRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	CountryCode string `json:"country_code" validate:"required,max=10"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,strong_password"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type BrandResponse struct {
	Id     int              `json:"id"`
	Logo   string           `json:"logo"`
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
	Password   string `json:"password" validate:"required,strong_password"`
	RePassword string `json:"re_password" validate:"required,strong_password,password_match"`
}

type ProfileResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	ImageUrl    string    `json:"image_url"`
	CountryCode string    `json:"country_code"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProfileEditRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	CountryCode string `json:"country_code" validate:"required,max=10"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type CarCreateRequest struct {
	ID                string    `json:"-"`
	OwnerId           string    `json:"-"`
	Title             string    `json:"title" validate:"required"`
	Description       string    `json:"description" validate:"required"`
	Currency          string    `json:"currency" validate:"required,currency"`
	Price             float64   `json:"price" validate:"required"`
	City              string    `json:"city" validate:"required"`
	District          string    `json:"district" validate:"required"`
	Neighborhood      string    `json:"neighborhood" validate:"required"`
	ListingNumber     string    `json:"-"`
	ListingDate       time.Time `json:"-"`
	BrandId           int       `json:"brand_id" validate:"required"`
	SeriesId          int       `json:"series_id" validate:"required"`
	ModelId           int       `json:"model_id" validate:"required"`
	Year              int       `json:"year" validate:"required"`
	FuelType          string    `json:"fuel_type" validate:"required,fuel_type"`
	Transmission      string    `json:"transmission" validate:"required,transmission"`
	Mileage           int       `json:"mileage" validate:"required"`
	BodyType          string    `json:"body_type" validate:"required,body_type"`
	EnginePower       int       `json:"engine_power" validate:"required"`
	EngineVolume      int       `json:"engine_volume" validate:"required"`
	DriveType         string    `json:"drive_type" validate:"required,drive_type"`
	Color             string    `json:"color" validate:"required"`
	Warranty          bool      `json:"warranty" validate:"required"`
	HeavyDamage       bool      `json:"heavy_damage"`
	SellerType        string    `json:"seller_type" validate:"required,seller_type"`
	TradeOption       bool      `json:"trade_option"`
	FrontBumper       string    `json:"front_bumper" validate:"required,condition"`
	FrontHood         string    `json:"front_hood" validate:"required,condition"`
	Roof              string    `json:"roof" validate:"required,condition"`
	FrontRightDoor    string    `json:"front_right_door" validate:"required,condition"`
	RearRightDoor     string    `json:"rear_right_door" validate:"required,condition"`
	FrontLeftMudguard string    `json:"front_left_mudguard" validate:"required,condition"`
	FrontLeftDoor     string    `json:"front_left_door" validate:"required,condition"`
	RearLeftDoor      string    `json:"rear_left_door" validate:"required,condition"`
	RearLeftMudguard  string    `json:"rear_left_mudguard" validate:"required,condition"`
	RearBumper        string    `json:"rear_bumper" validate:"required,condition"`
	Images            []string  `json:"images" validate:"omitempty"`
}

func (r CarCreateRequest) ToCar() *Car {
	return &Car{
		ID:                r.ID,
		OwnerId:           r.OwnerId,
		Title:             r.Title,
		Description:       r.Description,
		Currency:          r.Currency,
		Price:             r.Price,
		City:              r.City,
		District:          r.District,
		Neighborhood:      r.Neighborhood,
		ListingNumber:     r.ListingNumber,
		ListingDate:       r.ListingDate,
		BrandId:           r.BrandId,
		SeriesId:          r.SeriesId,
		ModelId:           r.ModelId,
		Year:              r.Year,
		FuelType:          r.FuelType,
		Transmission:      r.Transmission,
		Mileage:           r.Mileage,
		BodyType:          r.BodyType,
		EnginePower:       r.EnginePower,
		EngineVolume:      r.EngineVolume,
		DriveType:         r.DriveType,
		Color:             r.Color,
		Warranty:          r.Warranty,
		HeavyDamage:       r.HeavyDamage,
		SellerType:        r.SellerType,
		TradeOption:       r.TradeOption,
		FrontBumper:       r.FrontBumper,
		FrontHood:         r.FrontHood,
		Roof:              r.Roof,
		FrontRightDoor:    r.FrontRightDoor,
		RearRightDoor:     r.RearRightDoor,
		FrontLeftMudguard: r.FrontLeftMudguard,
		FrontLeftDoor:     r.FrontLeftDoor,
		RearLeftDoor:      r.RearLeftDoor,
		RearLeftMudguard:  r.RearLeftMudguard,
		RearBumper:        r.RearBumper,
	}
}

type ListCarResponse struct {
	Id          string    `json:"id,omitempty"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	Currency    string    `json:"currency"`
	Price       float64   `json:"price"`
	Brand       string    `json:"brand,omitempty"`
	Series      string    `json:"series,omitempty"`
	Model       string    `json:"model,omitempty"`
	Title       string    `json:"title,omitempty"`
	Year        int       `json:"year,omitempty"`
	Mileage     int       `json:"mileage,omitempty"`
	ListingDate time.Time `json:"listing_date,omitempty"`
	City        string    `json:"city,omitempty"`
	District    string    `json:"district,omitempty"`
}

type CarDetailResponse struct {
	ID                string        `json:"id,omitempty"`
	Owner             OwnerResponse `json:"owner,omitempty"`
	Title             string        `json:"title,omitempty"`
	Description       string        `json:"description,omitempty"`
	Currency          string        `json:"currency,omitempty"`
	Price             float64       `json:"price,omitempty"`
	City              string        `json:"city,omitempty"`
	District          string        `json:"district,omitempty"`
	Neighborhood      string        `json:"neighborhood,omitempty"`
	ListingNumber     string        `json:"listing_number,omitempty"`
	ListingDate       time.Time     `json:"listing_date,omitempty"`
	Brand             string        `json:"brand,omitempty"`
	Series            string        `json:"series,omitempty"`
	Model             string        `json:"model,omitempty"`
	Year              int           `json:"year,omitempty"`
	FuelType          string        `json:"fuel_type,omitempty"`
	Transmission      string        `json:"transmission,omitempty"`
	Mileage           int           `json:"mileage,omitempty"`
	BodyType          string        `json:"body_type,omitempty"`
	EnginePower       int           `json:"engine_power,omitempty"`
	EngineVolume      int           `json:"engine_volume,omitempty"`
	DriveType         string        `json:"drive_type,omitempty"`
	Color             string        `json:"color,omitempty"`
	Warranty          bool          `json:"warranty,omitempty"`
	HeavyDamage       bool          `json:"heavy_damage,omitempty"`
	SellerType        string        `json:"seller_type,omitempty"`
	TradeOption       bool          `json:"trade_option,omitempty"`
	FrontBumper       string        `json:"front_bumper,omitempty"`
	FrontHood         string        `json:"front_hood,omitempty"`
	Roof              string        `json:"roof,omitempty"`
	FrontRightDoor    string        `json:"front_right_door,omitempty"`
	RearRightDoor     string        `json:"rear_right_door,omitempty"`
	FrontLeftMudguard string        `json:"front_left_mudguard,omitempty"`
	FrontLeftDoor     string        `json:"front_left_door,omitempty"`
	RearLeftDoor      string        `json:"rear_left_door,omitempty"`
	RearLeftMudguard  string        `json:"rear_left_mudguard,omitempty"`
	RearBumper        string        `json:"rear_bumper,omitempty"`
	Images            []string      `json:"images,omitempty"`
}

type OwnerResponse struct {
	Id          string    `json:"id,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	CountryCode string    `json:"country_code,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
