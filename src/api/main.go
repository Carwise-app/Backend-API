package main

import (
	"carwise"
	"infra"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}

	app := gin.Default()
	app.Static("/images", "./images")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	interactor = carwise.NewInteractor(
		carwise.Services{
			UserRepo:          infra.NewUserRepository(),
			TokenRepo:         infra.NewTokenRepository(),
			MailGW:            infra.NewMailGateway(),
			PasswordResetRepo: infra.NewPasswordResetRepository(),
			CDNRepo:           infra.NewCDNRepository(),
			MessageRepo:       infra.NewMessageRepository(),
		},
	)

	auth := app.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
		auth.POST("/logout", Logout)
		auth.POST("/reset-password", ForgotPassword)
		auth.PUT("/reset-password", ResetPassword)
	}

	profile := app.Group("/profile")
	{
		profile.GET("/", AuthMiddleware(), Profile)
		profile.PUT("/edit", AuthMiddleware(), ProfileEdit)
	}
	/*
		aux := app.Group("/aux")
		{
			aux.GET("/brands", getBrands)
		}

		cars := app.Group("/cars")
		{
			cars.GET("/", listCars)
			cars.GET("/:id", getCarByID)
			cars.POST("/", AuthMiddleware(), createCar)
			cars.PUT("/:id", AuthMiddleware(), updateCar)
			cars.DELETE("/:id", AuthMiddleware(), deleteCar)
		}

		model := app.Group("/model")
		{
			model.POST("/predicts", predictPrice)
			model.GET("/predicts/history", AuthMiddleware(), getPredictionHistory)
			model.POST("/suggestions", suggestCar)
			model.GET("/suggestions/history", AuthMiddleware(), getSuggestionHistory)
		}
	*/

	app.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}
