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
			UserRepo: infra.NewUserRepository()},
	)

	auth := app.Group("/auth")
	{
		auth.POST("/register", registerUser)
		auth.POST("/login", loginUser)
		auth.POST("/logout", logoutUser)
		auth.POST("/reset-password", resetPasswordRequest)
		auth.PUT("/reset-password", resetPassword)
		auth.PUT("/edit", editUserProfile)
	}

	aux := app.Group("/aux")
	{
		aux.GET("/brands", getBrands)
		aux.GET("/brands/models", getModelsByBrand)
	}

	cars := app.Group("/cars")
	{
		cars.GET("/feed", getCarsFeed)
		cars.GET("/", listCars)
		cars.GET("/:id", getCarByID)
		cars.POST("/", createCar)
		cars.PUT("/:id", updateCar)
		cars.DELETE("/:id", deleteCar)
	}

	model := app.Group("/model")
	{
		model.POST("/predicts", predictPrice)
		model.GET("/predicts/history", getPredictionHistory)
		model.POST("/suggestions", suggestCar)
		model.GET("/suggestions/history", getSuggestionHistory)
	}

	app.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}
