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
	JWT_ISSUER string = "Carwise API Server"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

type UserClaims struct {
	UserId      string `json:"user_id"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
	Status      int    `json:"status"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

func JWTAuthorization(user *carwise.User) (string, error) {
	expirationTime := time.Now().Add(24 * 365 * time.Hour).Unix()
	claims := UserClaims{
		UserId:      user.Id,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		CountryCode: user.CountryCode,
		PhoneNumber: user.PhoneNumber,
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
		if !ok || !token.Valid || claims.Status == 2 || claims.Status == 3 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		ctx.Set("user", claims)
		ctx.Set("token", tokenString)
		ctx.Next()
	}
}

func WebSocketAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenString string

		// 1. Query parameter'dan token almaya çalış (?token=...)
		tokenString = ctx.Query("token")

		// 2. Eğer query'de yoksa Authorization header'dan al
		if tokenString == "" {
			authHeader := ctx.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 3. Eğer hala yoksa Sec-WebSocket-Protocol header'dan al (alternatif yöntem)
		if tokenString == "" {
			protocol := ctx.GetHeader("Sec-WebSocket-Protocol")
			if protocol != "" {
				// Protocol header'da "token.JWT_TOKEN_HERE" formatında gelebilir
				if strings.HasPrefix(protocol, "token.") {
					tokenString = strings.TrimPrefix(protocol, "token.")
				}
			}
		}

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token required for WebSocket connection"})
			ctx.Abort()
			return
		}

		// Token blacklist kontrolü
		isBlacklisted, errorMessages := interactor.IsTokenBlackListed(tokenString)
		if errorMessages == nil && isBlacklisted {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			ctx.Abort()
			return
		}

		// Token doğrulama
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
		if !ok || !token.Valid || claims.Status == 2 || claims.Status == 3 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		// WebSocket context'ine user bilgilerini ekle
		ctx.Set("user", claims)
		ctx.Set("token", tokenString)
		ctx.Next()
	}
}
