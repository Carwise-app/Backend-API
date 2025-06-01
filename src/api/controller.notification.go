package main

import (
	"carwise"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Push Notification
// @Description Push Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body carwise.PushNotificationRequest true "Push Notification Request"
// @Success 200 {object} map[string]interface{} "Push Notification Response" example:{"message":"Push notification sent successfully"}
// @Failure 401 {object} map[string]interface{} "Unauthorized" example:{"error":"No User found in request context"}
// @Failure 400 {object} map[string]interface{} "Bad Request" example:{"error":"Invalid request"}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"error":"Internal server error"}
// @Router /notification/push [post]
func PushNotification(c *gin.Context) {
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.PushNotificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	err := interactor.PushNotificationToAll(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Push notification sent successfully"})
}
