package carwise

type Listing struct {
	Id                string
	BrandId           string
	SeriesId          string
	ModelId           string
	Slug              string
	Title             string
	Description       string
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

type Image struct {
	Id   string
	Path string
}
