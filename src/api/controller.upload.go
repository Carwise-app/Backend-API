package main

import (
	"carwise"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Upload an image
// @Description Upload an image file (max 5MB, jpeg/png/gif)
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload"
// @Security BearerAuth
// @Success 200 {object} carwise.Image "Successfully uploaded image"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /upload/ [post]
func UploadImage(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.UploadImageRequest
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.File = file
	request.UserId = claim.UserId
	request.Role = claim.Role

	image, err := interactor.UploadImage(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"image": image})
}

// @Summary Delete an image
// @Description Delete an image by its ID
// @Tags Upload
// @Produce json
// @Param id path string true "Image ID"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse "Successfully deleted image"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /upload/{id} [delete]
func DeleteImage(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image id is required"})
		return
	}

	var request carwise.DeleteImageRequest
	request.ImageId = id
	request.UserId = claim.UserId
	request.Role = claim.Role

	err := interactor.DeleteImage(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "image deleted successfully"})
}

// @Summary Predict car damage from image
// @Description Use AI to predict if an image shows car damage
// @Tags Upload
// @Produce json
// @Param id path string true "Image ID"
// @Security BearerAuth
// @Success 200 {object} carwise.PredictImageResponse "Prediction results"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /upload/{id}/predict [get]
func PredictImage(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.PredictImageRequest
	request.ImageId = ctx.Param("id")
	request.UserId = claim.UserId
	request.Role = claim.Role

	response, err := interactor.PredictImage(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"prediction": response})
}

// @model ErrorResponse
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// @model SuccessResponse
// @Description Success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
