package main

import (
	"carwise"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenResponse represents the JWT token response
// @Description JWT token response
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// @Summary Verify Google ID Token
// @Description Verify Google ID token and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param idToken query string true "Google ID Token"
// @Success 200 {object} carwise.TokenResponse "Returns access token"
// @Failure 400 {object} map[string]interface{} "Invalid token or validation error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/google/id-token [get]
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

	responseToken := carwise.TokenResponse{
		AccessToken: token,
	}

	c.JSON(http.StatusOK, responseToken)
}

// @Summary Initiate Google Login
// @Description Redirect to Google OAuth login page
// @Tags auth
// @Produce json
// @Success 302 "Redirect to Google login page"
// @Router /auth/google/login [get]
func GoogleLogin(c *gin.Context) {
	url := interactor.GoogleAuthUrl()
	c.Redirect(http.StatusFound, url)
}

// @Summary Google OAuth Callback
// @Description Handle Google OAuth callback and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param state query string true "OAuth state parameter"
// @Param code query string true "OAuth authorization code"
// @Success 200 {object} carwise.TokenResponse "Returns access token"
// @Failure 400 {object} map[string]interface{} "Invalid state or code"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /auth/google/callback [get]
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
	responseToken := carwise.TokenResponse{
		AccessToken: token,
	}

	c.JSON(http.StatusOK, responseToken)
}
