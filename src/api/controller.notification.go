package main

import (
	"carwise"
	"log"
	"net/http"
	"strconv"

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

// @Summary Get Notifications
// @Description Get Notifications
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Success 200 {object} carwise.GetNotificationsResponse "Get Notifications Response"
// @Failure 401 {object} map[string]interface{} "Unauthorized" example:{"error":"No User found in request context"}
// @Failure 400 {object} map[string]interface{} "Bad Request" example:{"error":"Invalid request"}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"error":"Internal server error"}
// @Router /notification [get]
func GetNotifications(c *gin.Context) {
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.GetNotificationsRequest

	request.UserId = claim.UserId

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		log.Println("Invalid limit", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	request.Limit = limit

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Println("Invalid page", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}
	request.Page = page

	response, err := interactor.GetNotifications(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Read Notification
// @Description Read Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Read Notification Response" example:{"message":"Notification read successfully"}
// @Failure 401 {object} map[string]interface{} "Unauthorized" example:{"error":"No User found in request context"}
// @Failure 400 {object} map[string]interface{} "Bad Request" example:{"error":"Invalid request"}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"error":"Internal server error"}
// @Router /notification/{notificationId} [put]
func ReadNotification(c *gin.Context) {
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.ReadNotificationRequest
	request.UserId = claim.UserId

	notificationId := c.Param("id")
	request.NotificationId = notificationId

	err := interactor.ReadNotification(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification read successfully"})
}

// @Summary Delete Notification
// @Description Delete Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Delete Notification Response" example:{"message":"Notification deleted successfully"}
// @Failure 401 {object} map[string]interface{} "Unauthorized" example:{"error":"No User found in request context"}
// @Failure 400 {object} map[string]interface{} "Bad Request" example:{"error":"Invalid request"}
// @Failure 500 {object} map[string]interface{} "Internal Server Error" example:{"error":"Internal server error"}
// @Router /notification/{notificationId} [delete]
func DeleteNotification(c *gin.Context) {
	userContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.DeleteNotificationRequest
	request.UserId = claim.UserId

	notificationId := c.Param("id")
	request.NotificationId = notificationId

	err := interactor.DeleteNotification(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
