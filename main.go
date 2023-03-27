package main

//go:generate oapi-codegen --package=api -generate=types -o ./api/storygenie.gen.go https://api.swaggerhub.com/apis/swagger354/storygenie/0.0.2

import (
	"log"
	"os"
	"storygenie-backend/controller"
	"storygenie-backend/database"
	"storygenie-backend/helper"
	"storygenie-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	loadEnv()
	helper.GetFirebaseApp()
	serveApplication()
}

func loadDatabase() *gorm.DB {
	database := database.Connect()
	// database.AutoMigrate(&models.Story{})
	// database.AutoMigrate(&models.Product{})
	// database.AutoMigrate(&models.Feedback{})
	return database
}

func loadEnv() {
	godotenv.Load(".env")
}

func serveApplication() {
	pCtrl := controller.PublicController{Database: loadDatabase()}
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3200", "https://app.storygenie.io"}
	router.Use(cors.New(config))
	router.GET("/health", pCtrl.HealthCheck)
	privateRoutes := router.Group("/api")
	privateRoutes.Use(middleware.Authentication)
	// if os.Getenv("ENVIRONMENT") != "production" {
	// 	privateRoutes.GET("/seed", pCtrl.SeedDatabase)
	// }
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

	port := os.Getenv("PORT")
	if port == "" {
		log.Print("Defaulting to port 8000")
		port = "8000"
	}

	log.Print("Hello from Cloud Run! The container started successfully and is listening for HTTP requests on $PORT")
	log.Printf("Listening on port %s", port)
	router.Run(":" + port)
}
