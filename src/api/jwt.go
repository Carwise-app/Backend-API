package main

import (
	"carwise"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	JWT_ISSUER string = "CARWISE API SERVER"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

type UserClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
	jwt.StandardClaims
}

func JWTAuthorization(user *carwise.User) (string, error) {
	expirationTime := time.Now().Add(24 * 365 * time.Hour).Unix()
	claims := UserClaims{
		UserId: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
			Issuer:    JWT_ISSUER,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		isBlacklisted, errorMessages := interactor.IsTokenBlackListed(tokenString)
		if errorMessages == nil && isBlacklisted {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, gin.Error{
					Err:  http.ErrAbortHandler,
					Type: gin.ErrorTypePrivate,
				}
			}
			return JWT_SECRET, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		ctx.Set("user", claims)
		ctx.Set("token", tokenString)
		ctx.Next()
	}
}
