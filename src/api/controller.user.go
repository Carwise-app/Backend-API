package main

import (
	"carwise"
	"log"
	"net/http"
	"strconv"

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
// @Router /auth/reset-password [post]
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
// @Router /auth/reset-password [put]
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

// @Summary Get user by ID
// @Description Get user by ID
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} carwise.UserInfo "User information"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /user/{id} [get]
func GetUserById(ctx *gin.Context) {
	var request carwise.GetUserByIdRequest

	request.Id = ctx.Param("id")
	user, errors := interactor.GetUserById(request)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Count
// @Description Count
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} carwise.CountResponse "Count response"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /admin/count [get]
func Count(ctx *gin.Context) {
	var request carwise.CountRequest

	claim := GetUserClaims(ctx)

	request.UserId = claim.UserId
	request.Role = claim.Role

	count, errors := interactor.Count(&request)
	if errors != nil {
		log.Println("Error counting users:", errors.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.JSON(http.StatusOK, count)
}

// @Summary Get users
// @Description Get users
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int true "Page number"
// @Param limit query int true "Limit number"
// @Success 200 {object} carwise.GetUsersResponse "Get users response"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /admin/users [get]
func GetUsers(ctx *gin.Context) {
	var request carwise.GetUsersRequest

	claim := GetUserClaims(ctx)

	request.UserId = claim.UserId
	request.Role = claim.Role

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		log.Println("Error parsing page parameter:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	request.Page = page

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		log.Println("Error parsing limit parameter:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	request.Limit = limit

	users, errors := interactor.GetUsers(&request)
	if errors != nil {
		log.Println("Error getting users:", errors.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// @Summary Delete user
// @Description Delete user
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /profile/{id} [delete]
func DeleteUser(ctx *gin.Context) {
	var request carwise.DeleteAccountRequest

	claim := GetUserClaims(ctx)

	request.UserId = claim.UserId
	request.Role = claim.Role
	request.ProfileId = ctx.Param("id")

	errors := interactor.DeleteAccount(request)
	if errors != nil {
		log.Println("Error deleting user:", errors.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Update user role
// @Description Update user role (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body carwise.UpdateUserRoleRequest true "Update user role request"
// @Success 200 {object} map[string]interface{} "User role updated successfully" example:{"message":"User role updated successfully"}
// @Failure 400 {object} map[string]interface{} "Validation error" example:{"error":"Invalid role value"}
// @Failure 401 {object} map[string]interface{} "Unauthorized" example:{"error":"No User found in request context"}
// @Failure 403 {object} map[string]interface{} "Forbidden" example:{"error":"Only admin users can update user roles"}
// @Failure 500 {object} map[string]interface{} "Server error" example:{"error":"Internal server error"}
// @Router /admin/update-user-role [put]
func UpdateUserRole(ctx *gin.Context) {
	claim := GetUserClaims(ctx)

	var request carwise.UpdateUserRoleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println("Error binding JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set admin information from JWT claims
	request.AdminUserId = claim.UserId
	request.AdminRole = claim.Role

	// Validate request structure
	errors := ValidateStruct(request)
	if errors != nil {
		log.Println("Validation errors:", errors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	// Call interactor to update user role
	errors = interactor.UpdateUserRole(request)
	if errors != nil {
		log.Println("Error updating user role:", errors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}
