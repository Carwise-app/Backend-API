package main

import (
	"carwise"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Register a new user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body carwise.UserCreateRequest true "User registration request"
// @Success 200 {object} carwise.TokenResponse "Returns access token"
// @Failure 400 {object} map[string]interface{} "Validation or creation error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/register [post]
func Register(ctx *gin.Context) {
	var request carwise.UserCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println("Error binding JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errors := ValidateStruct(request)
	if errors != nil {
		log.Println("Validation errors:", errors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	user, errors := interactor.CreateUser(request)
	if errors != nil {
		log.Println("Error creating user:", errors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	token, err := JWTAuthorization(user)
	if err != nil {
		log.Println("Error generating token:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": []string{"Could not generate token"}})
		return
	}

	response := carwise.TokenResponse{
		AccessToken: token,
	}

	ctx.JSON(http.StatusOK, response)

}

// @Summary Login user
// @Description Login user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body carwise.UserLoginRequest true "User login request"
// @Success 200 {object} carwise.TokenResponse "Returns access token"
// @Failure 400 {object} map[string]interface{} "Validation or login error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/login [post]
func Login(ctx *gin.Context) {
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

	response := carwise.TokenResponse{
		AccessToken: token,
	}

	ctx.JSON(http.StatusOK, response)

}

// @Summary Logout user
// @Description Logout user by blacklisting their token
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 "Successfully logged out"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/logout [post]
func Logout(ctx *gin.Context) {
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

// @Summary Request password reset
// @Description Send password reset email to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body carwise.ResetPasswordRequest true "Password reset request"
// @Success 200 "Reset email sent successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Router /auth/forgot-password [post]
func ForgotPassword(ctx *gin.Context) {
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

// @Summary Reset password
// @Description Reset user password using token and email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Reset token"
// @Param email query string true "User email"
// @Param request body carwise.ChangePasswordRequest true "New password request"
// @Success 200 "Password reset successful"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Router /auth/reset-password [post]
func ResetPassword(ctx *gin.Context) {
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
