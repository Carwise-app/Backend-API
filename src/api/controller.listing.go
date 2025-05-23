package main

import (
	"carwise"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateListing(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.CreateListingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	id, err := interactor.CreateListing(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Listing created successfully", "id": id})
}

func GetListing(ctx *gin.Context) {
	id := ctx.Param("id")

	listing, err := interactor.GetListingById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, listing)
}

func ListListing(ctx *gin.Context) {
	var request carwise.ListListingRequest

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	request.Filter.Page = pageInt
	request.Filter.Limit = limitInt

	query := ctx.Query("query")
	if query != "" {
		request.Filter.Query = query
	}

	bodyType := ctx.Query("body_type")
	if bodyType != "" {
		request.Filter.BodyType = bodyType
	}

	driveType := ctx.Query("drive_type")
	if driveType != "" {
		request.Filter.DriveType = driveType
	}

	transmissionType := ctx.Query("transmission_type")
	if transmissionType != "" {
		request.Filter.TransmissionType = transmissionType
	}

	fuelType := ctx.Query("fuel_type")
	if fuelType != "" {
		request.Filter.FuelType = fuelType
	}

	city := ctx.Query("city")
	if city != "" {
		request.Filter.City = city
	}

	district := ctx.Query("district")
	if district != "" {
		request.Filter.District = district
	}

	neighborhood := ctx.Query("neighborhood")
	if neighborhood != "" {
		request.Filter.Neighborhood = neighborhood
	}

	brandId := ctx.Query("brand_id")
	if brandId != "" {
		request.Filter.BrandId = brandId
	}

	seriesId := ctx.Query("series_id")
	if seriesId != "" {
		request.Filter.SeriesId = seriesId
	}

	modelId := ctx.Query("model_id")
	if modelId != "" {
		request.Filter.ModelId = modelId
	}

	minPrice := ctx.Query("min_price")
	if minPrice != "" {
		minPriceInt, err := strconv.Atoi(minPrice)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min price"})
			return
		}
		request.Filter.MinPrice = minPriceInt
	}
	maxPrice := ctx.Query("max_price")
	if maxPrice != "" {
		maxPriceInt, err := strconv.Atoi(maxPrice)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max price"})
			return
		}
		request.Filter.MaxPrice = maxPriceInt
	}

	minYear := ctx.Query("min_year")
	if minYear != "" {
		minYearInt, err := strconv.Atoi(minYear)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min year"})
			return
		}
		request.Filter.MinYear = minYearInt
	}
	maxYear := ctx.Query("max_year")
	if maxYear != "" {
		maxYearInt, err := strconv.Atoi(maxYear)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max year"})
			return
		}
		request.Filter.MaxYear = maxYearInt
	}

	minKilometers := ctx.Query("min_kilometers")
	if minKilometers != "" {
		minKilometersInt, err := strconv.Atoi(minKilometers)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min kilometers"})
			return
		}
		request.Filter.MinKilometers = minKilometersInt
	}
	maxKilometers := ctx.Query("max_kilometers")
	if maxKilometers != "" {
		maxKilometersInt, err := strconv.Atoi(maxKilometers)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max kilometers"})
			return
		}
		request.Filter.MaxKilometers = maxKilometersInt
	}

	minEnginePower := ctx.Query("min_engine_power")
	if minEnginePower != "" {
		minEnginePowerInt, err := strconv.Atoi(minEnginePower)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min engine power"})
			return
		}
		request.Filter.MinEnginePower = minEnginePowerInt
	}
	maxEnginePower := ctx.Query("max_engine_power")
	if maxEnginePower != "" {
		maxEnginePowerInt, err := strconv.Atoi(maxEnginePower)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max engine power"})
			return
		}
		request.Filter.MaxEnginePower = maxEnginePowerInt
	}

	minEngineVolume := ctx.Query("min_engine_volume")
	if minEngineVolume != "" {
		minEngineVolumeInt, err := strconv.Atoi(minEngineVolume)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min engine volume"})
			return
		}
		request.Filter.MinEngineVolume = minEngineVolumeInt
	}
	maxEngineVolume := ctx.Query("max_engine_volume")
	if maxEngineVolume != "" {
		maxEngineVolumeInt, err := strconv.Atoi(maxEngineVolume)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max engine volume"})
			return
		}
		request.Filter.MaxEngineVolume = maxEngineVolumeInt
	}

	color := ctx.Query("color")
	if color != "" {
		request.Filter.Color = color
	}

	heavyDamage := ctx.Query("heavy_damage")
	if heavyDamage != "" {
		heavyDamageBool, err := strconv.ParseBool(heavyDamage)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid heavy damage"})
			return
		}
		request.Filter.HeavyDamage = &heavyDamageBool
	}

	createdBy := ctx.Query("created_by")
	if createdBy != "" {
		request.Filter.CreatedBy = createdBy
	}

	request.Filter.SortBy = ctx.DefaultQuery("sort", "created_at")
	if !slices.Contains([]string{"created_at", "price", "kilometers", "year"}, request.Filter.SortBy) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort field"})
		return
	}
	request.Filter.SortOrder = strings.ToLower(ctx.DefaultQuery("order", "asc"))
	if !slices.Contains([]string{"asc", "desc"}, request.Filter.SortOrder) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort order"})
		return
	}

	status := ctx.Query("status")
	if status != "" {
		statusInt, err := strconv.Atoi(status)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		request.Filter.Status = statusInt
	}

	response, err := interactor.ListListing(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func UpdateListing(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.UpdateListingRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing id"})
		return
	}
	request.Id = id

	err := interactor.UpdateListing(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Listing updated successfully"})
}

func DeleteListing(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.DeleteListingRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing id"})
		return
	}
	request.Id = id

	err := interactor.DeleteListing(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

func UpdateListingStatus(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.UpdateListingStatusRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing id"})
		return
	}
	request.Id = id

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = interactor.UpdateListingStatus(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Listing status updated successfully"})
}
