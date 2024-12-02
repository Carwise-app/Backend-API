package main

import (
	"carwise"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func userProfile(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	profile, err := interactor.GetProfile(claim.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

func resetPasswordRequest(ctx *gin.Context) {
	var request carwise.ResetPasswordRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": []string{err.Error()},
		})
		return
	}

	if err := ValidateStruct(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if errors := interactor.ResetPasswordRequest(request); errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.Status(http.StatusOK)
}
func resetPassword(ctx *gin.Context) {
	var request carwise.ChangePasswordRequest

	token := ctx.Query("token")
	email := ctx.Query("email")

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": []string{err.Error()},
		})
		return
	}

	if err := ValidateStruct(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if errors := interactor.ChangePassword(request, token, email); errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.Status(http.StatusOK)

}

func editUserProfile(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.ProfileEditRequest

	request.FirstName = ctx.Request.FormValue("first_name")
	request.LastName = ctx.Request.FormValue("last_name")
	request.CountryCode = ctx.Request.FormValue("country_code")
	request.PhoneNumber = ctx.Request.FormValue("phone_number")

	if err := ValidateStruct(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	avatar, err := ctx.FormFile("avatar")
	if avatar != nil && err != nil {
		if err.Error() != "multipart: no multipart data" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": []string{err.Error()},
			})
			return
		}
	}

	if avatar != nil && !isValidImageFormat(avatar.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": []string{"Invalid file format."},
		})
		return
	}

	if errors := interactor.EditProfile(claim.UserId, request, avatar); errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func getBrands(ctx *gin.Context) {
	brands, err := interactor.GetBrands()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, brands)
}

func listCars(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "20")
	pageStr := ctx.DefaultQuery("page", "1")
	brandIDStr := ctx.DefaultQuery("brand_id", "0")
	seriesIDStr := ctx.DefaultQuery("series_id", "0")
	modelIDStr := ctx.DefaultQuery("model_id", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println("Error converting limit:", err)
		ctx.JSON(400, gin.H{"error": "Invalid limit"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Println("Error converting page:", err)
		ctx.JSON(400, gin.H{"error": "Invalid page"})
		return
	}

	brandID, err := strconv.Atoi(brandIDStr)
	if err != nil {
		log.Println("Error converting brand_id:", err)
		ctx.JSON(400, gin.H{"error": "Invalid brand_id"})
		return
	}

	seriesID, err := strconv.Atoi(seriesIDStr)
	if err != nil {
		log.Println("Error converting series_id:", err)
		ctx.JSON(400, gin.H{"error": "Invalid series_id"})
		return
	}

	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		log.Println("Error converting model_id:", err)
		ctx.JSON(400, gin.H{"error": "Invalid model_id"})
		return
	}

	response, errors := interactor.ListCars(page, limit, brandID, seriesID, modelID)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.JSON(http.StatusOK, response)

}
func getCarByID(ctx *gin.Context) {
	id := ctx.Param("id")

	carDetail, err := interactor.GetCarDetail(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, carDetail)

}
func createCar(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.CarCreateRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": []string{err.Error()},
		})
		return
	}

	if err := ValidateStruct(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	if errors := interactor.CreateCar(claim.UserId, request); errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.Status(http.StatusOK)

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

func isValidImageFormat(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}
