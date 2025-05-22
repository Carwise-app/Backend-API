package carwise

type Brand struct {
	Id      string
	ImageId string
	Name    string
}

type Series struct {
	Id      string
	BrandId string
	Name    string
}

type Model struct {
	Id       int
	BrandId  string
	SeriesId string
	Name     string
}
