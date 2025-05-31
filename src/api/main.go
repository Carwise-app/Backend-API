// Package main is the entry point for the Carwise API server
// @title Carwise API
// @version 1.0
// @description Carwise API documentation
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"carwise"
	_ "docs" // Import the generated docs package with blank identifier
	"infra"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Carwise API
// @version         1.0
// @description     Carwise API documentation with Swagger
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.email  support@carwise.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      carwisegw.yusuftalhaklc.com
// @BasePath  /

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}

	app := gin.Default()
	app.Static("/uploads", "./uploads")
	app.Static("/.well-known", "./well-known")

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
			RedisRepo:         infra.NewRedisRepository(),
			MailGW:            infra.NewMailGateway(),
			PasswordResetRepo: infra.NewPasswordResetRepository(),
			BrandRepo:         infra.NewBrandRepository(),
			ListingRepo:       infra.NewListingRepository(),
			ImageRepo:         infra.NewImageRepository(),
			MessageRepo:       infra.NewMessageRepository(),
			PredictionRepo:    infra.NewPredictionRepository(),
			GoogleAuth:        infra.NewGoogleAuth(),
			FavoriteRepo:      infra.NewFavoriteRepository(),
			PredictRepo:       infra.NewPredictRepository(),
		},
	)

	go hub.Run()

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	app.GET("/user/:id", AuthMiddleware(), GetUserById)

	auth := app.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
		auth.POST("/logout", Logout)
		auth.POST("/reset-password", ForgotPassword)
		auth.PUT("/reset-password", ResetPassword)
		auth.GET("/google", GoogleLogin)
		auth.GET("/google/callback", GoogleCallback)
		auth.GET("/google/id-token", GoogleIdToken)
		auth.GET("/reset-password", func(c *gin.Context) {
			token := c.Query("token")
			email := c.Query("email")
			c.Redirect(http.StatusMovedPermanently, "carwise://reset-password?token="+token+"&email="+email)
		})
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
		listing.GET("/:id", OptionalAuthMiddleware(), GetListing)
		listing.PUT("/:id", AuthMiddleware(), UpdateListing)
		listing.DELETE("/:id", AuthMiddleware(), DeleteListing)
		listing.PATCH("/:id/status", AuthMiddleware(), UpdateListingStatus)
		listing.GET("/", OptionalAuthMiddleware(), ListListing)
	}

	upload := app.Group("/upload")
	{
		upload.POST("", AuthMiddleware(), UploadImage)
		upload.DELETE("/:id", AuthMiddleware(), DeleteImage)
		upload.GET("/:id/predict", AuthMiddleware(), PredictImage)
	}

	chat := app.Group("/chat")
	{
		chat.GET("/", AuthMiddleware(), GetChats)
		chat.POST("/:listing_id/:receiver_id", AuthMiddleware(), SendMessage)
		chat.GET("/:listing_id/:receiver_id", AuthMiddleware(), GetMessages)
		chat.GET("/ws", WebSocketAuthMiddleware(), WebSocketHandler)
	}

	favorite := app.Group("/favorite")
	{
		favorite.POST("/:listing_id", AuthMiddleware(), CreateFavorite)
		favorite.DELETE("/:listing_id", AuthMiddleware(), DeleteFavorite)
		favorite.GET("/", AuthMiddleware(), GetFavorites)
	}

	predict := app.Group("/predict")
	{
		predict.POST("/", OptionalAuthMiddleware(), CreatePredict)
		predict.GET("/", AuthMiddleware(), GetPredicts)
	}

	app.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))
	openBrowser(getSwaggerURL())
	app.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		return
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}

	if err != nil {
		log.Printf("Failed to open browser: %v\n", err)
	}
}

func getSwaggerURL() string {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	swaggerURL := "http://" + host + ":" + port + "/swagger/index.html"
	return swaggerURL
}
