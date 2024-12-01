package carwise

import (
	"time"
)

const (
	FuelTypeDiesel       = "Diesel"
	FuelTypePetrol       = "Petrol"
	FuelTypePetrolAndLPG = "Petrol & LPG"
	FuelTypeLPG          = "LPG"
	FuelTypeElectric     = "Electric"
	//-----------------------------
	TransmissionAutomatic = "Automatic"
	TransmissionManual    = "Manual"
	//-----------------------------
	BodyTypeSedan     = "Sedan"
	BodyTypeHatchback = "Hatchback"
	BodyTypeCoupe     = "Coupe"
	//----------------------------
	PartConditionOriginal = "Original"
	PartConditionPainted  = "Painted"
	PartConditionChanged  = "Changed"
	//------------------------------
	SellerTypeIndividual = "Individual"
	SellerTypeDealer     = "Dealer"
	//--------------------------------
	CurrencyTRY = "TRY"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
	//------------------------------
	DriveTypeFrontWheelDrive = "Front-Wheel Drive"
	DriveTypeRearWheelDrive  = "Rear-Wheel Drive"
	DriveTypeFourWheelDrive  = "Four-Wheel Drive"
	DriveTypeAllWheelDrive   = "All-Wheel Drive"
)

type Car struct {
	ID                string
	Title             string
	Description       string
	Currency          string
	Price             float64
	City              string
	District          string
	Neighborhood      string
	ListingNumber     string
	ListingDate       time.Time
	BrandId           int
	SeriesId          int
	ModelId           int
	Year              int
	FuelType          string
	Transmission      string
	Condition         string
	Mileage           int
	BodyType          string
	EnginePower       int
	EngineVolume      int
	DriveType         string
	Color             string
	Warranty          bool
	HeavyDamage       bool
	PlateCountry      string
	SellerType        string
	TradeOption       bool
	FrontBumper       string
	FrontHood         string
	Roof              string
	FrontRightDoor    string
	RearRightDoor     string
	FrontLeftMudguard string
	FrontLeftDoor     string
	RearLeftDoor      string
	RearLeftMudguard  string
	RearBumper        string
}

type Images struct {
	ID    int
	CarId string
	URL   string
}

type Brand struct {
	ID   int
	Logo string
	Name string
}

type Series struct {
	ID      int
	BrandID string
	Name    string
}

type Model struct {
	ID       int
	SeriesID string
	Name     string
}

const (
	UserRoleAdmin   = "Admin"
	UserRoleRegular = "Regular"
	//------------------------
	AccountStatusActive   = "Active"
	AccountStatusInactive = "Inactive"
	AccountStatusBanned   = "Banned"
)

type User struct {
	ID           string
	FirstName    string
	LastName     string
	ImageUrl     string
	CountryCode  string
	PhoneNumber  string
	Email        string
	PasswordHash string
	Role         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLogin    time.Time
}
