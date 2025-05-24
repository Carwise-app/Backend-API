package carwise

// @model Listing
// @Description Listing information
type Listing struct {
	Id                string
	BrandId           string
	SeriesId          string
	ModelId           string
	Slug              string
	Title             string
	Description       string
	Status            int
	Currency          string
	Price             int
	City              string
	District          string
	Neighborhood      string
	Images            []string
	FuelType          string
	TransmissionType  string
	BodyType          string
	DriveType         string
	EnginePower       int
	EngineVolume      int
	Kilometers        int
	Year              int
	Color             string
	HeavyDamage       bool
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
	CreatedBy         string
	CreatedAt         int64
	UpdatedAt         int64
}

// @model Image
// @Description Image information
type Image struct {
	Id        string `json:"id" example:"img_123456"`            // Unique identifier of the image
	Path      string `json:"path" example:"/uploads/img123.jpg"` // Path where the image is stored
	CreatedBy string `json:"created_by" example:"user_123"`      // ID of the user who uploaded the image
	CreatedAt int64  `json:"created_at" example:"1646092800"`    // Timestamp when the image was uploaded
}

type ImagePrediction struct {
	ImageId    string  `json:"image_id"`
	Prediction bool    `json:"prediction"`
	Confidence float64 `json:"confidence"`
	CreatedAt  int64   `json:"created_at"`
}
