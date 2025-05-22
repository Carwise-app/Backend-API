package carwise

type CreateListingRequest struct {
	BrandId      string            `json:"brand_id"`
	SeriesId     string            `json:"series_id"`
	ModelId      string            `json:"model_id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Currency     string            `json:"currency"`
	Price        int               `json:"price"`
	City         string            `json:"city"`
	District     string            `json:"district"`
	Neighborhood string            `json:"neighborhood"`
	Images       []string          `json:"images"`
	DetailInfo   ListingDetailInfo `json:"detail"`
	UserId       string            `json:"-"`
	Role         int               `json:"-"`
}

type ListingDetailInfo struct {
	FuelType          string `json:"fuel_type"`
	TransmissionType  string `json:"transmission_type"`
	BodyType          string `json:"body_type"`
	DriveType         string `json:"drive_type"`
	EnginePower       int    `json:"engine_power"`
	EngineVolume      int    `json:"engine_volume"`
	Kilometers        int    `json:"kilometers"`
	Year              int    `json:"year"`
	Color             string `json:"color"`
	HeavyDamage       bool   `json:"heavy_damage"`
	FrontBumper       string `json:"front_bumper"`
	FrontHood         string `json:"front_hood"`
	Roof              string `json:"roof"`
	FrontRightDoor    string `json:"front_right_door"`
	RearRightDoor     string `json:"rear_right_door"`
	FrontLeftMudguard string `json:"front_left_mudguard"`
	FrontLeftDoor     string `json:"front_left_door"`
	RearLeftDoor      string `json:"rear_left_door"`
	RearLeftMudguard  string `json:"rear_left_mudguard"`
	RearBumper        string `json:"rear_bumper"`
}

type GetListingResponse struct {
	Id           string            `json:"id"`
	Slug         string            `json:"slug"`
	Brand        Brand             `json:"brand"`
	Series       Series            `json:"series"`
	Model        Model             `json:"model"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Currency     string            `json:"currency"`
	Price        int               `json:"price"`
	City         string            `json:"city"`
	District     string            `json:"district"`
	Neighborhood string            `json:"neighborhood"`
	Images       []Image           `json:"images"`
	DetailInfo   ListingDetailInfo `json:"detail"`
	CreatedBy    UserInfo          `json:"created_by"`
	CreatedAt    int64             `json:"created_at"`
	UpdatedAt    int64             `json:"updated_at"`
}
