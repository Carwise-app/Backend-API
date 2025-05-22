package main

import (
	"carwise"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Profile(ctx *gin.Context) {
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

func ProfileEdit(ctx *gin.Context) {
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

func isValidImageFormat(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}
