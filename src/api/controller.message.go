package main

import (
	"carwise"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SendMessage(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.SendMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.UserId = claim.UserId
	request.Role = claim.Role

	request.ReceiverId = ctx.Param("receiver_id")

	err := interactor.SendMessage(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func GetMessages(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.GetMessagesRequest

	request.ReceiverId = ctx.Param("receiver_id")
	request.Limit, _ = strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	request.Page, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))

	request.UserId = claim.UserId
	request.Role = claim.Role

	response, err := interactor.GetMessages(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func GetChats(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.GetChatsRequest
	request.Limit, _ = strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	request.Page, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	request.UserId = claim.UserId
	request.Role = claim.Role

	response, err := interactor.GetChats(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
