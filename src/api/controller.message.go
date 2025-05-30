package main

import (
	"carwise"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Send a message
// @Description Send a message to a user
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param receiver_id path string true "Receiver ID"
// @Param listing_id path string true "Listing ID"
// @Param message body carwise.SendMessageRequest true "Message details"
// @Success 200 {object} map[string]interface{} "Message sent successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /chat/{receiver_id} [post]
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
	request.ListingId = ctx.Param("listing_id")

	err := interactor.SendMessage(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// @Summary Get messages
// @Description Get messages between two users
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param receiver_id path string true "Receiver ID"
// @Param listing_id path string true "Listing ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} carwise.GetMessagesResponse "Messages retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /chat/{receiver_id} [get]
func GetMessages(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	var request carwise.GetMessagesRequest

	request.ReceiverId = ctx.Param("receiver_id")
	request.ListingId = ctx.Param("listing_id")
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

// @Summary Get chats
// @Description Get all chats for a user
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} carwise.GetChatsResponse "Chats retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /chat [get]
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

func WebSocketHandler(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.String(http.StatusUnauthorized, "No User found in request context")
		return
	}
	claim := userContext.(*UserClaims)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserId: claim.UserId,
		Role:   claim.Role,
	}

	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
