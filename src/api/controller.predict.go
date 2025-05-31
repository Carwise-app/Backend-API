package main

import (
	"carwise"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new car price prediction
// @Description Create a new car price prediction using the ML model
// @Tags Predict
// @Accept json
// @Produce json
// @Param request body carwise.PredictRequest true "Car details for prediction"
// @Security BearerAuth
// @Success 200 {object} carwise.PredictionResponse "Prediction results"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /predict [post]
func CreatePredict(c *gin.Context) {
	var request carwise.PredictRequest
	userContext, exists := c.Get("user")
	if exists {
		claim := userContext.(*UserClaims)
		request.UserId = claim.UserId
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	predicts, err := interactor.CreatePredict(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, predicts)
}

// @Summary Get user's predictions
// @Description Get a paginated list of predictions made by the user
// @Tags Predict
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Security BearerAuth
// @Success 200 {object} carwise.GetPredictsResponse "List of predictions"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /predict [get]
func GetPredicts(c *gin.Context) {
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.GetPredictsRequest
	request.UserId = claim.UserId

	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	request.Page = pageInt
	request.Limit = limitInt

	predicts, err := interactor.GetPredicts(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, predicts)
}
