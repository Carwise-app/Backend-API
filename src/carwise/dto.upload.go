package carwise

import "mime/multipart"

type UploadImageRequest struct {
	File   *multipart.FileHeader `json:"file"`
	UserId string                `json:"-"`
	Role   int                   `json:"-"`
}

type DeleteImageRequest struct {
	ImageId string `json:"image_id"`
	UserId  string `json:"-"`
	Role    int    `json:"-"`
}
