package carwise

type BrandCreateRequest struct {
	ImageId string `json:"image_id"`
	Name    string `json:"name" validate:"required"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
}

type BrandUpdateRequest struct {
	ImageId string `json:"image_id"`
	Name    string `json:"name" validate:"required"`
	BrandId string `json:"-"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
}

type BrandDeleteRequest struct {
	BrandId string `json:"-"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
}

type SeriesCreateRequest struct {
	BrandId string `json:"-"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
	Name    string `json:"name" validate:"required"`
}

type SeriesUpdateRequest struct {
	SeriesId string `json:"-"`
	BrandId  string `json:"-"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
	Name     string `json:"name" validate:"required"`
}

type SeriesDeleteRequest struct {
	SeriesId string `json:"-"`
	BrandId  string `json:"-"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}

type ModelCreateRequest struct {
	BrandId  string `json:"-"`
	SeriesId string `json:"-"`
	Name     string `json:"name" validate:"required"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}

type ModelUpdateRequest struct {
	ModelId  string `json:"-"`
	BrandId  string `json:"-"`
	SeriesId string `json:"-"`
	Name     string `json:"name" validate:"required"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}

type ModelDeleteRequest struct {
	ModelId  string `json:"-"`
	BrandId  string `json:"-"`
	SeriesId string `json:"-"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}
