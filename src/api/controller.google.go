package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoogleIdToken(c *gin.Context) {
	idToken := c.Query("idToken")
	response, err := interactor.GoogleIdToken(idToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := JWTAuthorization(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func GoogleLogin(c *gin.Context) {
	url := interactor.GoogleAuthUrl()
	c.Redirect(http.StatusFound, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state not found"})
		return
	}
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}
	account, err := interactor.GoogleAuthCallback(state, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := JWTAuthorization(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"access_token": token,
	})
}
