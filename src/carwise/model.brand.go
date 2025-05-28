package carwise

type Brand struct {
	Id        string `json:"id,omitempty"`
	ImagePath string `json:"image_path,omitempty"`
	Name      string `json:"name,omitempty"`
}

type Series struct {
	Id      string `json:"id,omitempty"`
	BrandId string `json:"brand_id,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Model struct {
	Id       string `json:"id,omitempty"`
	BrandId  string `json:"brand_id,omitempty"`
	SeriesId string `json:"series_id,omitempty"`
	Name     string `json:"name,omitempty"`
}

type BrandWithDetails struct {
	Id        string         `json:"id,omitempty"`
	ImagePath string         `json:"image_path,omitempty"`
	Name      string         `json:"name,omitempty"`
	Series    []SeriesDetail `json:"series,omitempty"`
}

type SeriesDetail struct {
	Id     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Models []Model `json:"models,omitempty"`
}
