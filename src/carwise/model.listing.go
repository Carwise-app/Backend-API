package carwise

type Listing struct {
	Id           string
	BrandId      string
	SeriesId     string
	ModelId      string
	Title        string
	Description  string
	Currency     string
	Price        int
	City         string
	District     string
	Neighborhood string
	Images       []string
	CreatedBy    string
	CreatedAt    int64
	UpdatedAt    int64
}

type Image struct {
	Id   string
	Path string
}
