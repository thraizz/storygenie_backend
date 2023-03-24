package controller

import (
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicController struct {
	Database *gorm.DB
}

func (c *PublicController) SeedDatabase(context *gin.Context) {
	c.Database.AutoMigrate(&models.User{})
	c.Database.AutoMigrate(&models.Product{})
	c.Database.AutoMigrate(&models.Prompt{})
	c.Database.AutoMigrate(&models.Story{})
}
