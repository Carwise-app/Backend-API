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
