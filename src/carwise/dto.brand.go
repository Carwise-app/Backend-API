package carwise

// @model BrandCreateRequest
// @Description Request body for creating a new brand
type BrandCreateRequest struct {
	ImagePath string `json:"image_path" example:"img_123456"`
	Name      string `json:"name" validate:"required" example:"BMW"`
	UserId    string `json:"-"`
	Role      int    `json:"-"`
}

// @model BrandUpdateRequest
// @Description Request body for updating an existing brand
type BrandUpdateRequest struct {
	ImagePath string `json:"image_path" example:"img_123456"`
	Name      string `json:"name" validate:"required" example:"BMW"`
	BrandId   string `json:"-"`
	UserId    string `json:"-"`
	Role      int    `json:"-"`
}

// @model BrandDeleteRequest
// @Description Request body for deleting a brand
type BrandDeleteRequest struct {
	BrandId string `json:"-"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
}

// @model SeriesCreateRequest
// @Description Request body for creating a new series
type SeriesCreateRequest struct {
	BrandId string `json:"-"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
	Name    string `json:"name" validate:"required" example:"3 Series"`
}

// @model SeriesUpdateRequest
// @Description Request body for updating an existing series
type SeriesUpdateRequest struct {
	SeriesId string `json:"-"`
	BrandId  string `json:"-"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
	Name     string `json:"name" validate:"required" example:"3 Series"`
}

// @model SeriesDeleteRequest
// @Description Request body for deleting a series
type SeriesDeleteRequest struct {
	SeriesId string `json:"-"`
	BrandId  string `json:"-"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}

// @model ModelCreateRequest
// @Description Request body for creating a new model
type ModelCreateRequest struct {
	BrandId  string `json:"-"`
	SeriesId string `json:"-"`
	Name     string `json:"name" validate:"required" example:"320i"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}

// @model ModelUpdateRequest
// @Description Request body for updating an existing model
type ModelUpdateRequest struct {
	ModelId  string `json:"-"`
	BrandId  string `json:"-"`
	SeriesId string `json:"-"`
	Name     string `json:"name" validate:"required" example:"320i"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}

// @model ModelDeleteRequest
// @Description Request body for deleting a model
type ModelDeleteRequest struct {
	ModelId  string `json:"-"`
	BrandId  string `json:"-"`
	SeriesId string `json:"-"`
	UserId   string `json:"-"`
	Role     int    `json:"-"`
}
