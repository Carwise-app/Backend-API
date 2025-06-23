package main

import (
	"carwise"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create a favorite
// @Description Create a favorite for a listing
// @Tags Favorite
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param listing_id path string true "Listing ID"
// @Success 201 {object} map[string]interface{} "Favorite created successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /favorite/{listing_id} [post]
func CreateFavorite(ctx *gin.Context) {
	claim := GetUserClaims(ctx)

	var request carwise.FavoriteRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	listingId := ctx.Param("listing_id")
	request.ListingId = listingId

	err := interactor.CreateFavorite(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Favorite created successfully"})
}

// @Summary Delete a favorite
// @Description Delete a favorite for a listing
// @Tags Favorite
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param listing_id path string true "Listing ID"
// @Success 200 {object} map[string]interface{} "Favorite deleted successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /favorite/{listing_id} [delete]
func DeleteFavorite(ctx *gin.Context) {
	claim := GetUserClaims(ctx)

	var request carwise.DeleteFavoriteRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	listingId := ctx.Param("listing_id")
	request.ListingId = listingId

	err := interactor.DeleteFavorite(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// @Summary Get favorites
// @Description Get favorites for a user
// @Tags Favorite
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} carwise.ListListingResponse "Favorites retrieved successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /favorite [get]
func GetFavorites(ctx *gin.Context) {
	claim := GetUserClaims(ctx)

	var request carwise.GetFavoritesRequest
	request.UserId = claim.UserId
	request.Role = claim.Role

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	request.Limit = limit
	request.Page = page

	response, err := interactor.GetFavorites(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
