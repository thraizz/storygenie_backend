package main

//go:generate oapi-codegen --package=api -generate=types -o ./api/storygenie.gen.go https://api.swaggerhub.com/apis/swagger354/storygenie/0.0.3

import (
	"fmt"
	"log"
	"os"
	"storygenie-backend/controller"
	"storygenie-backend/database"
	"storygenie-backend/helper"
	"storygenie-backend/middleware"
	"storygenie-backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
)

func main() {
	loadEnv()
	helper.GetFirebaseApp()
	serveApplication()
}

func initializeSentry() {
	SENTRY_DSN_KEY := os.Getenv("SENTRY_DSN_KEY")
	if SENTRY_DSN_KEY == "" {
		panic("Sentry DSN key not found")
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              SENTRY_DSN_KEY,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	log.Print("Sentry initialized.")
}

func loadDatabase() *gorm.DB {
	database := database.Connect()
	database.AutoMigrate(&models.Story{})
	database.AutoMigrate(&models.Prompt{})
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.Feedback{})
	database.AutoMigrate(&models.JiraRefreshToken{})
	return database
}

func loadEnv() {
	godotenv.Load(".env")
}

func serveApplication() {
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	initializeSentry()
	pCtrl := controller.PublicController{Database: loadDatabase()}
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3200", "https://app.storygenie.io"}
	router.Use(cors.New(config))
	router.GET("/health", pCtrl.HealthCheck)
	privateRoutes := router.Group("/api")
	privateRoutes.Use(middleware.Authentication)
	privateRoutes.GET("/jira/secret", pCtrl.GetJiraClientSecret)
	privateRoutes.GET("/jira/refresh", pCtrl.GetJiraRefreshToken)
	privateRoutes.PUT("/jira/refresh", pCtrl.SetJiraRefreshToken)
	privateRoutes.GET("/story", pCtrl.GetStories)
	privateRoutes.GET("/story/:storyId", pCtrl.GetStoryById)
	privateRoutes.POST("/story", pCtrl.CreateStory)
	privateRoutes.DELETE("/story/:storyId", pCtrl.DeleteStory)
	privateRoutes.GET("/story/:storyId/feedback", pCtrl.GetFeedbackForStory)
	privateRoutes.POST("/story/:storyId/feedback", pCtrl.CreateFeedback)
	privateRoutes.POST("/story/generate", pCtrl.GenerateScrumStories)
	privateRoutes.GET("/product", pCtrl.GetProducts)
	privateRoutes.GET("/product/:productId", pCtrl.GetProductById)
	privateRoutes.POST("/product", pCtrl.CreateProduct)
	privateRoutes.DELETE("/product/:productId", pCtrl.DeleteProduct)
	if os.Getenv("ENVIRONMENT") == "dev" {
		privateRoutes.GET("/seed", pCtrl.SeedDatabase)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Print("No PORT env var, defaulting to port 8000")
		port = "8000"
	}

	log.Printf("Hello from %s The container started successfully and is listening for HTTP requests on %s", os.Getenv("ENVIRONMENT"), port)
	router.Run(":" + port)
}
