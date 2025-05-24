package main

import (
	"carwise"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all brands with details
// @Description Retrieve all brands with their associated details
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Success 200 {object} []carwise.BrandWithDetails
// @Failure 500 {object} map[string]string
// @Router /brands [get]
func GetAllBrands(ctx *gin.Context) {
	brands, err := interactor.GetAllBrandsWithDetails()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"brands": brands})
}

// @Summary Create a new brand
// @Description Create a new brand with the provided details
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param brand body carwise.BrandCreateRequest true "Brand creation request"
// @Success 200 {object} map[string]string{message=string} "Brand created successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /brands [post]
func CreateBrand(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.BrandCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	err := interactor.CreateBrand(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Brand created successfully"})
}

// @Summary Update an existing brand
// @Description Update a brand with the provided details
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param brand body carwise.BrandUpdateRequest true "Brand update request"
// @Success 200 {object} map[string]string{message=string} "Brand updated successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /brands/{id} [put]
func UpdateBrand(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.BrandUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	err := interactor.UpdateBrand(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Brand updated successfully"})
}

// @Summary Delete a brand
// @Description Delete a brand by its ID
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Success 200 {object} map[string]string{message=string} "Brand deleted successfully"
// @Failure 401 {object} map[string]string
// @Router /brands/{id} [delete]
func DeleteBrand(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.BrandDeleteRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	err := interactor.DeleteBrand(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Brand deleted successfully"})
}

// @Summary Create a new series
// @Description Create a new series for a specific brand
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param series body carwise.SeriesCreateRequest true "Series creation request"
// @Success 200 {object} map[string]string{message=string} "Series created successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /brands/{id}/series [post]
func CreateSeries(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.SeriesCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	err := interactor.CreateSeries(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Series created successfully"})
}

// @Summary Update an existing series
// @Description Update a series with the provided details
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param sid path string true "Series ID"
// @Param series body carwise.SeriesUpdateRequest true "Series update request"
// @Success 200 {object} map[string]string{message=string} "Series updated successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /brands/{id}/series/{sid} [put]
func UpdateSeries(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.SeriesUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	seriesId := ctx.Param("sid")
	request.SeriesId = seriesId

	err := interactor.UpdateSeries(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Series updated successfully"})
}

// @Summary Delete a series
// @Description Delete a series by its ID
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param sid path string true "Series ID"
// @Success 200 {object} map[string]string{message=string} "Series deleted successfully"
// @Failure 401 {object} map[string]string
// @Router /brands/{id}/series/{sid} [delete]
func DeleteSeries(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.SeriesDeleteRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	seriesId := ctx.Param("sid")
	request.SeriesId = seriesId

	err := interactor.DeleteSeries(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Series deleted successfully"})
}

// @Summary Create a new model
// @Description Create a new model for a specific series
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param sid path string true "Series ID"
// @Param model body carwise.ModelCreateRequest true "Model creation request"
// @Success 200 {object} map[string]string{message=string} "Model created successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /brands/{id}/series/{sid}/models [post]
func CreateModel(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.ModelCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	seriesId := ctx.Param("sid")
	request.SeriesId = seriesId

	err := interactor.CreateModel(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Model created successfully"})
}

// @Summary Update an existing model
// @Description Update a model with the provided details
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param sid path string true "Series ID"
// @Param mid path string true "Model ID"
// @Param model body carwise.ModelUpdateRequest true "Model update request"
// @Success 200 {object} map[string]string{message=string} "Model updated successfully"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /brands/{id}/series/{sid}/models/{mid} [put]
func UpdateModel(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.ModelUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	seriesId := ctx.Param("sid")
	request.SeriesId = seriesId

	err := interactor.UpdateModel(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Model updated successfully"})
}

// @Summary Delete a model
// @Description Delete a model by its ID
// @Tags Brands - Series - Models
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param sid path string true "Series ID"
// @Param mid path string true "Model ID"
// @Success 200 {object} map[string]string{message=string} "Model deleted successfully"
// @Failure 401 {object} map[string]string
// @Router /brands/{id}/series/{sid}/models/{mid} [delete]
func DeleteModel(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.ModelDeleteRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	brandId := ctx.Param("id")
	request.BrandId = brandId

	seriesId := ctx.Param("sid")
	request.SeriesId = seriesId

	modelId := ctx.Param("mid")
	request.ModelId = modelId

	err := interactor.DeleteModel(&request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Model deleted successfully"})
}
