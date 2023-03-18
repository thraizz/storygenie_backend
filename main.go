package main

import (
	"fmt"
	"log"
	"storygenie-backend/controller"
	"storygenie-backend/database"
	"storygenie-backend/helper"
	"storygenie-backend/middleware"
	"storygenie-backend/models"

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
	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Project{})
	database.AutoMigrate(&models.Prompt{})
	database.AutoMigrate(&models.Story{})
	return database
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func serveApplication() {
	pCtrl := controller.PublicController{Database: loadDatabase()}
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	router.GET("/seed", pCtrl.SeedDatabase)
	privateRoutes := router.Group("/api")
	privateRoutes.Use(middleware.Authentication)
	privateRoutes.GET("/story", pCtrl.GetStories)
	privateRoutes.GET("/story/:storyId", pCtrl.GetStoryById)
	privateRoutes.POST("/story", pCtrl.CreateStory)
	privateRoutes.GET("/project", pCtrl.GetProjects)
	privateRoutes.GET("/project/:storyId", pCtrl.GetProjects)
	privateRoutes.POST("/project", pCtrl.CreateProject)
	privateRoutes.DELETE("/project/:projectId", pCtrl.DeleteProject)

	router.Run(":8000")
	fmt.Println("Server running on port 8000")
}
