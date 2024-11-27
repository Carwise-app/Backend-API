package main

import (
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

	auth := app.Group("/auth")
	{
		auth.POST("/register")
		auth.POST("/login")
		auth.POST("/logout")
		auth.POST("/reset-password")
		auth.PUT("/reset-password")
		auth.PUT("/edit")
	}

	aux := app.Group("/aux")
	{
		aux.GET("/brands")
		aux.GET("/brands/models")
	}

	cars := app.Group("/cars")
	{
		cars.GET("/feed")
		cars.GET("/")
		cars.GET("/:id")
		cars.POST("/")
		cars.PUT("/:id")
		cars.DELETE("/:id")
	}

	model := app.Group("/model")
	{
		model.POST("/predicts")
		model.GET("/predicts/history")
		model.POST("/suggestions")
		model.GET("/suggestions/history")
	}

	app.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}
