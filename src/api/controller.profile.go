package main

import (
	"carwise"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Get user profile
// @Description Get the profile information of the authenticated user
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} carwise.ProfileResponse "User profile information"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /profile [get]
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

// @Summary Edit user profile
// @Description Update the profile information of the authenticated user
// @Tags Profile
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param first_name formData string true "First name" example:"John"
// @Param last_name formData string true "Last name" example:"Doe"
// @Param country_code formData string true "Country code" example:"+90"
// @Param phone_number formData string true "Phone number" example:"5551234567"
// @Param avatar formData file false "Profile picture"
// @Success 200 "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /profile [put]
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

// @Summary Notify user profile
// @Description Notify the profile information of the authenticated user
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Param notify_request body carwise.ProfileNotifyRequest true "Notify request"
// @Success 200 "Profile notify updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /profile/notify [patch]
func ProfileNotify(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.ProfileNotifyRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.UserId = claim.UserId

	if err := interactor.NotifyProfile(request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
