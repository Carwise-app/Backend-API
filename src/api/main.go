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
	app.Static("/uploads", "./uploads")

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
			BrandRepo:         infra.NewBrandRepository(),
			ListingRepo:       infra.NewListingRepository(),
			ImageRepo:         infra.NewImageRepository(),
			MessageRepo:       infra.NewMessageRepository(),
		},
	)

	// WebSocket hub'ı başlat
	go hub.Run()

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

	brand := app.Group("/brand")
	{
		brand.POST("/", AuthMiddleware(), CreateBrand)
		brand.PUT("/:id", AuthMiddleware(), UpdateBrand)
		brand.DELETE("/:id", AuthMiddleware(), DeleteBrand)
		brand.GET("/", AuthMiddleware(), GetAllBrands)

		brand.POST("/:id/series", AuthMiddleware(), CreateSeries)
		brand.PUT("/:id/series/:sid", AuthMiddleware(), UpdateSeries)
		brand.DELETE("/:id/series/:sid", AuthMiddleware(), DeleteSeries)

		brand.POST("/:id/series/:sid/model", AuthMiddleware(), CreateModel)
		brand.PUT("/:id/series/:sid/model/:mid", AuthMiddleware(), UpdateModel)
		brand.DELETE("/:id/series/:sid/model/:mid", AuthMiddleware(), DeleteModel)
	}

	listing := app.Group("/listing")
	{
		listing.POST("/", AuthMiddleware(), CreateListing)
		listing.GET("/:id", GetListing)
		listing.PUT("/:id", AuthMiddleware(), UpdateListing)
		listing.DELETE("/:id", AuthMiddleware(), DeleteListing)
		listing.PATCH("/:id/status", AuthMiddleware(), UpdateListingStatus)
		listing.GET("/", ListListing)
	}

	upload := app.Group("/upload")
	{
		upload.POST("/", AuthMiddleware(), UploadImage)
		upload.DELETE("/:id", AuthMiddleware(), DeleteImage)
	}

	chat := app.Group("/chat")
	{
		chat.GET("/", AuthMiddleware(), GetChats)
		chat.POST("/:receiver_id", AuthMiddleware(), SendMessage)
		chat.GET("/:receiver_id", AuthMiddleware(), GetMessages)
		chat.GET("/ws", WebSocketAuthMiddleware(), WebSocketHandler)
	}

	app.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}
