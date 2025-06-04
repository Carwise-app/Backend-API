package main

import (
	"carwise"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new listing
// @Description Create a new listing with the given details
// @Tags Listing Car
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param listing body carwise.CreateListingRequest true "Listing details"
// @Success 201 {object} carwise.IdResponse "Listing created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /listing [post]
func CreateListing(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.CreateListingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	id, err := interactor.CreateListing(&request)
	if err != nil {
		fmt.Println(request)
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Get a listing by ID
// @Description Get a listing by its unique identifier
// @Tags Listing Car
// @Produce json
// @Param id path string true "Listing ID"
// @Success 200 {object} carwise.GetListingResponse "Listing details"
// @Failure 404 {object} ErrorResponse "Listing not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /listing/{id} [get]
func GetListing(ctx *gin.Context) {
	var request carwise.GetListingRequest
	userContext, exists := ctx.Get("user")
	if exists {
		claim := userContext.(*UserClaims)
		request.UserId = claim.UserId
	}
	id := ctx.Param("id")
	request.Id = id

	listing, err := interactor.GetListingById(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, listing)
}

// @Summary List listings
// @Description Get a list of listings with optional filters
// @Tags Listing
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param query query string false "Search query"
// @Param brand_id query string false "Brand ID"
// @Param series_id query string false "Series ID"
// @Param model_id query string false "Model ID"
// @Param body_type query string false "Body type"
// @Param drive_type query string false "Drive type"
// @Param transmission_type query string false "Transmission type"
// @Param fuel_type query string false "Fuel type"
// @Param city query string false "City"
// @Param district query string false "District"
// @Param neighborhood query string false "Neighborhood"
// @Param min_price query int false "Minimum price"
// @Param max_price query int false "Maximum price"
// @Param min_year query int false "Minimum year"
// @Param max_year query int false "Maximum year"
// @Param min_kilometers query int false "Minimum kilometers"
// @Param max_kilometers query int false "Maximum kilometers"
// @Param min_engine_power query int false "Minimum engine power"
// @Param max_engine_power query int false "Maximum engine power"
// @Param min_engine_volume query int false "Minimum engine volume"
// @Param max_engine_volume query int false "Maximum engine volume"
// @Param color query string false "Color"
// @Param heavy_damage query bool false "Heavy damage"
// @Param sort query string false "Sort field (created_at, price, kilometers, year)" default(created_at)
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Param status query int false "Listing status"
// @Param created_by query string false "Created by user ID"
// @Success 200 {object} carwise.ListListingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /listing [get]
func ListListing(ctx *gin.Context) {
	var request carwise.ListListingRequest

	// Get user context if available
	userContext, exists := ctx.Get("user")
	if exists {
		claim := userContext.(*UserClaims)
		request.Filter.UserId = claim.UserId
	}

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

// @Summary Update a listing
// @Description Update the details of a listing
// @Tags Listing Car
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Listing ID"
// @Param listing body carwise.UpdateListingRequest true "Listing details"
// @Success 200 {object} map[string]interface{} "Listing updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /listing/{id} [put]
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

// @Summary Delete a listing
// @Description Delete a listing by its unique identifier
// @Tags Listing Car
// @Security BearerAuth
// @Param id path string true "Listing ID"
// @Success 200 {object} map[string]interface{} "Listing deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /listing/{id} [delete]
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

// @Summary Update a listing status
// @Description Update the status of a listing
// @Tags Listing Car
// @Security BearerAuth
// @Param id path string true "Listing ID"
// @Param status body carwise.UpdateListingStatusRequest true "Listing status"
// @Success 200 {object} map[string]interface{} "Listing status updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /listing/{id}/status [patch]
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
