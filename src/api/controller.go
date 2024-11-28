package main

import (
	"carwise"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	interactor *carwise.Interactor
)

func registerUser(ctx *gin.Context) {
	var request carwise.UserCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errors := ValidateStruct(request)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	user, errors := interactor.CreateUser(request)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	token, err := JWTAuthorization(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": []string{"Could not generate token"}})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": token})

}
func loginUser(ctx *gin.Context) {
	var request carwise.UserLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errors := ValidateStruct(request)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	user, errors := interactor.LoginUser(request)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	token, err := JWTAuthorization(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": []string{"Could not generate token"}})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": token})

}

func logoutUser(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": []string{"No token found in request context"}})
		return
	}
	tokenString := token.(string)

	err := interactor.AddTokenBlackList(tokenString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.Status(http.StatusOK)
}

func resetPasswordRequest(c *gin.Context) {

}
func resetPassword(c *gin.Context) {

}
func editUserProfile(c *gin.Context) {

}

func getBrands(c *gin.Context) {

}
func getModelsByBrand(c *gin.Context) {

}

func getCarsFeed(c *gin.Context) {

}
func listCars(c *gin.Context) {

}
func getCarByID(c *gin.Context) {

}
func createCar(c *gin.Context) {

}
func updateCar(c *gin.Context) {

}
func deleteCar(c *gin.Context) {

}

func predictPrice(c *gin.Context) {

}
func getPredictionHistory(c *gin.Context) {

}
func suggestCar(c *gin.Context) {

}
func getSuggestionHistory(c *gin.Context) {

}
