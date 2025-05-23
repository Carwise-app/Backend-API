package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadImage(ctx *gin.Context) {
	userContext, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No User found in request context"})
		return
	}
	claim := userContext.(*UserClaims)

	ctx.JSON(http.StatusOK, gin.H{"image": claim.Email})
}
