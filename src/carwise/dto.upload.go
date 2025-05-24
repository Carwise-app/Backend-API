// Package carwise provides the data transfer objects and business logic for the Carwise application
package carwise

import "mime/multipart"

// @model UploadImageRequest
// @Description Request body for uploading an image
type UploadImageRequest struct {
	File   *multipart.FileHeader `json:"file" swaggerignore:"true"` // File to upload (max 5MB, jpeg/png/gif)
	UserId string                `json:"-"`                         // User ID (from context)
	Role   int                   `json:"-"`                         // User role (from context)
}

// @model DeleteImageRequest
// @Description Request body for deleting an image
type DeleteImageRequest struct {
	ImageId string `json:"image_id" example:"img_123456"` // ID of the image to delete
	UserId  string `json:"-"`                             // User ID (from context)
	Role    int    `json:"-"`                             // User role (from context)
}

// @model PredictImageRequest
// @Description Request body for predicting car damage from an image
type PredictImageRequest struct {
	ImageId string `json:"image_id" example:"img_123456"` // ID of the image to predict
	UserId  string `json:"-"`                             // User ID (from context)
	Role    int    `json:"-"`                             // User role (from context)
}

// @model PredictImageResponse
// @Description Response containing the prediction results for an image
type PredictImageResponse struct {
	Image      Image   `json:"image"`                           // Image information
	Prediction bool    `json:"prediction" example:"true"`       // Whether the image shows car damage
	Confidence float64 `json:"confidence" example:"0.95"`       // Confidence score of the prediction
	CreatedAt  int64   `json:"created_at" example:"1646092800"` // Timestamp of the prediction
}
