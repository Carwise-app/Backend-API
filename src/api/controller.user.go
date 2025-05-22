package main

import (
	"carwise"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{"access_token": token})

}

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
