package carwise

// @model CreateListingRequest
// @Description Request body for creating a new listing
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

// @model ListingDetailInfo
// @Description Detailed information about the listing
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

// @model GetListingResponse
// @Description Response body for getting a listing
type GetListingResponse struct {
	Id           string            `json:"id"`
	Slug         string            `json:"slug"`
	Status       int               `json:"status"`
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
	IsFavorite   bool              `json:"is_favorite"`
	Images       []Image           `json:"images"`
	DetailInfo   ListingDetailInfo `json:"detail"`
	CreatedBy    UserInfo          `json:"created_by"`
	CreatedAt    int64             `json:"created_at"`
	UpdatedAt    int64             `json:"updated_at"`
}

// @model ListListingRequest
// @Description Request body for listing listings
type ListListingRequest struct {
	Filter ListingFilter
}

// @model ListingFilter
// @Description Filter criteria for listing search
type ListingFilter struct {
	BrandId          string `json:"brand_id"`
	SeriesId         string `json:"series_id"`
	ModelId          string `json:"model_id"`
	Query            string `json:"query"`
	BodyType         string `json:"body_type"`
	DriveType        string `json:"drive_type"`
	TransmissionType string `json:"transmission_type"`
	FuelType         string `json:"fuel_type"`
	City             string `json:"city"`
	District         string `json:"district"`
	Neighborhood     string `json:"neighborhood"`
	MinPrice         int    `json:"min_price"`
	MaxPrice         int    `json:"max_price"`
	MinYear          int    `json:"min_year"`
	MaxYear          int    `json:"max_year"`
	MinKilometers    int    `json:"min_kilometers"`
	MaxKilometers    int    `json:"max_kilometers"`
	MinEnginePower   int    `json:"min_engine_power"`
	MaxEnginePower   int    `json:"max_engine_power"`
	MinEngineVolume  int    `json:"min_engine_volume"`
	MaxEngineVolume  int    `json:"max_engine_volume"`
	Color            string `json:"color"`
	HeavyDamage      *bool  `json:"heavy_damage"`
	SortBy           string `json:"sort_by"`
	SortOrder        string `json:"sort_order"`
	CreatedBy        string `json:"created_by"`
	Status           int    `json:"status"`
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	UserId           string `json:"-"`
}

// @model UpdateListingRequest
// @Description Request body for updating a listing
type UpdateListingRequest struct {
	Id           string            `json:"id"`
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

// @model DeleteListingRequest
// @Description Request body for deleting a listing
type DeleteListingRequest struct {
	Id     string `json:"-"`
	UserId string `json:"-"`
	Role   int    `json:"-"`
}

// @model UpdateListingStatusRequest
// @Description Request body for updating the status of a listing
type UpdateListingStatusRequest struct {
	Id     string `json:"-"`
	Status int    `json:"status" validate:"required,min=1,max=3"`
	UserId string `json:"-"`
	Role   int    `json:"-"`
}

// @model ListListingResponse
// @Description Response body for listing listings
type ListListingResponse struct {
	Listings []ListListingInfo `json:"listings"`
	Total    int               `json:"total"`
}

// @model ListListingInfo
// @Description Information about a listing
type ListListingInfo struct {
	Id           string `json:"id"`
	Slug         string `json:"slug"`
	Status       int    `json:"status"`
	Brand        Brand  `json:"brand"`
	Series       Series `json:"series"`
	Model        Model  `json:"model"`
	Title        string `json:"title"`
	Currency     string `json:"currency"`
	Price        int    `json:"price"`
	City         string `json:"city"`
	District     string `json:"district"`
	Neighborhood string `json:"neighborhood"`
	IsFavorite   bool   `json:"is_favorite"`
	Image        Image  `json:"image"`
	CreatedAt    int64  `json:"created_at"`
}

// @model IdResponse
// @Description Response body for id
type IdResponse struct {
	Id string `json:"id"`
}

// @model GetListingRequest
// @Description Request body for getting a listing
type GetListingRequest struct {
	Id     string `json:"id"`
	UserId string `json:"-"`
}
